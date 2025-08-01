package account

import (
	"log"
	"net"
	"strings"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	u "github.com/lemoras/goutils/api"
	d "github.com/lemoras/goutils/db"
)

type Membership struct {
	RoleId     int
	AppId      int
	MerchantId uuid.UUID
	HasId      bool
	CustomData string
}

// a struct to rep user account
type Account struct {
	gorm.Model
	//ID             int64     `gorm:"primaryKey"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	Token          string    `json:"token";sql:"-"`
	Nickname       string    `json:"nickname"`
	PhotoUrl       string    `json:"photoUrl"`
	AboutMe        string    `json:"aboutMe"`
	UserId         uuid.UUID `json:"userId"`
	LoginErrorLock string
	Lock           bool `json:"lock"`
}

// Validate incoming user details...
func (account *Account) Validate() (map[string]interface{}, bool) {

	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "0x11001:Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "0x11002:Password is required"), false
	}

	var domain = ""
	if strings.ContainsAny(account.Email, "@") {
		domain = strings.Split(account.Email, "@")[1]
	}

	var hasMX bool

	//check if domain has MX record
	mxRecords, errMx := net.LookupMX(domain)

	if errMx != nil {
		log.Println("No MX records found", errMx)
	}

	if len(mxRecords) > 0 {
		hasMX = true
	}

	if !hasMX {
		return u.Message(false, "0x11003:Email address is not valid"), false
	}

	//Email must be unique
	temp := &Account{}

	//check for errors and duplicate emails
	err := d.GetDB().Table("security.accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return u.Message(false, "0x11005:Email address already in use by another user"), false
	}

	return u.Message(true, "0x11006:Requirement passed"), true
}

func (account *Account) Create() map[string]interface{} {

	if resp, ok := account.Validate(); !ok {
		return resp
	}

	account.Lock = false
	account.LoginErrorLock = ""

	if account.Nickname == "" {
		account.Nickname = strings.Split(account.Email, "@")[0]
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)

	account.Password = string(hashedPassword)
	account.UserId = uuid.New()

	d.GetDB().Table("security.accounts").Create(account)

	if account.ID <= 0 {
		return u.Message(false, "0x11007:Failed to create account, connection error")
	}

	account.Token = ""

	account.Password = ""
	account.ID = 0

	response := u.Message(true, "0x11008:Account has been created")
	response["account"] = account
	return response
}

func (account *Account) Update(user uuid.UUID, hasId bool) map[string]interface{} {

	resp, ok, curAccount := GetUser(user)
	if !ok {
		return resp
	}

	isModify := false
	if account.Password != "" {

		if hasId {
			return u.Message(false, "0x11009:The password is changed only by its owner")
		}

		s := strings.Split(account.Password, "#change-password#")
		if len(s) != 2 {
			resp := u.Message(false, "0x11010:Password error")
			return resp
		}
		if len(s[1]) < 6 {
			resp := u.Message(false, "0x11002:Password is required")
			return resp
		}
		err := bcrypt.CompareHashAndPassword([]byte(curAccount.Password), []byte(s[0]))
		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
			return u.Message(false, "0x11011:Current Password is not valid")
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(s[1]), bcrypt.DefaultCost)
		curAccount.Password = string(hashedPassword)
		isModify = true
	}

	if account.Email != "" {
		if hasId {
			return u.Message(false, "0x11012:The email is changed only by its owner")
		}
		curAccount.Email = account.Email
		isModify = true
	}

	if account.Nickname != "" {
		curAccount.Nickname = account.Nickname
		isModify = true
	}
	if account.PhotoUrl != "" {
		curAccount.PhotoUrl = account.PhotoUrl
		isModify = true
	}
	if account.AboutMe != "" {
		curAccount.AboutMe = account.AboutMe
		isModify = true
	}

	if isModify {
		d.GetDB().Table("security.accounts").Save(curAccount)
	} else {
		return u.Message(false, "0x110013:Failed to modify account, can't update values")
	}

	curAccount.Password = "" //delete password
	curAccount.ID = 0

	response := u.Message(true, "0x11014:0x11047:Account has been modified")
	response["account"] = curAccount
	return response
}

func Get(uid uuid.UUID) map[string]interface{} {

	resp, ok, account := GetUser(uid)
	if !ok {
		return resp
	}

	account.ID = 0
	account.Password = ""
	response := u.Message(true, "0x11015:Account found it")
	response["account"] = account
	return response
}

func GetUser(userId uuid.UUID) (map[string]interface{}, bool, *Account) {

	acc := &Account{}
	err := d.GetDB().Table("security.accounts").Where("user_id = ?", userId).First(acc).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "0x11016:The user was not found"), false, acc
		}
		return u.Message(false, "0x11004:Connection error. Please retry"), false, acc
	}

	// if acc.UserId != userId {
	// 	return u.Message(false, "The user was not match given userId"), false, acc
	// }

	//acc.Password = ""    VERY IMPORTANT OTHER USE IT
	return u.Message(true, "0x11006:Requirement passed"), true, acc
}
