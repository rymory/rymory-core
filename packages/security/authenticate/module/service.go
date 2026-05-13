// Copyright (c) 2017-2026 Onur Yaşar
// Licensed under AGPL v3 + Commercial Exception
// See LICENSE.txt

package authenticate

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	u "github.com/lemoras/goutils/api"
	d "github.com/lemoras/goutils/db"
	"golang.org/x/crypto/bcrypt"
)

type Membership struct {
	UserId     uuid.UUID
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
	UserId         uuid.UUID `json:"-"`
	LoginErrorLock string    `json:"-"`
	Lock           bool      `json:"-"`
}

type BuildAccount struct {
	LastLoginDate time.Time `json:"lastLoginDate"`
	Token         string    `json:"token";sql:"-"`
	UserId        uuid.UUID `json:"userId"`
}

type Result struct {
	ID             int       `json:"-"`
	Email          string    `json:"email"`
	Password       string    `json:"-"`
	Token          string    `json:"token";sql:"-"`
	Nickname       string    `json:"nickname"`
	PhotoUrl       string    `json:"photoUrl"`
	AboutMe        string    `json:"aboutMe"`
	UserId         uuid.UUID `json:"-"`
	AppList        string    `json:"appList"`
	CreatedAt      time.Time `json:"createdAt"`
	LoginErrorLock string    `json:"-"`
	Lock           bool      `json:"-"`
}

type ResultBuildToken struct {
	UserId        uuid.UUID `json:"userId"`
	ProjectInfo   string    `json:"projectInfo"`
	LastLoginDate time.Time `json:"lastLoginDate"`
}

func Login(email, password string, tokenUserId uuid.UUID) map[string]interface{} {

	account := &Result{}

	var subject string

	err := d.GetDB().Raw("SELECT *, (SELECT string_agg(CONCAT (Membership.roles.app_id::TEXT, ':', Membership.roles.role_id::TEXT, ':', Membership.roles.merchant_id::TEXT, ':', Membership.roles.project_info::TEXT), ', ') AS app_list  FROM  Membership.roles  WHERE Membership.roles.user_id = security.accounts.user_id and Membership.roles.active and Membership.roles.deleted_at is null GROUP BY security.accounts.user_id) FROM security.accounts WHERE security.accounts.lock = ? and security.accounts.email = ?", false, email).Find(account).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "0x11021:Email address not found")

		}
		fmt.Println(err)
		return u.Message(false, "0x11004:Connection error. Please retry")
	}

	if len(account.LoginErrorLock) > 0 && (account.LoginErrorLock[0] == "3"[0] || account.LoginErrorLock[0] == "6"[0] || account.LoginErrorLock[0] == "9"[0]) {

		if account.LoginErrorLock[0] == "9"[0] {
			lockAccountResult := LockAccount(account.UserId)
			if !lockAccountResult["status"].(bool) {
				return lockAccountResult
			}
			return u.Message(false, "0x11138:Your account has been locked due to 9 unsuccessful login attempts")
		}

		lastTime, _ := strconv.ParseInt(strings.Split(account.LoginErrorLock, "x")[1], 10, 64)

		// deltaTime := time.Now().Unix() - time.Now().Add(-(time.Duration(300) * time.Second)).Unix()

		deltaTime := time.Now().Unix() - lastTime

		if deltaTime <= 300 && account.LoginErrorLock[0] == "3"[0] {
			return u.Message(false, "0x11139:You have entered incorrect credentials 3 times. You need to wait 5 minutes after your last unsuccessful attempt")
		}

		if deltaTime <= 900 && account.LoginErrorLock[0] == "6"[0] {
			return u.Message(false, "0x11140:You have entered incorrect credentials 6 times. You need to wait 15 minutes after your last unsuccessful attempt")
		}
	}

	if tokenUserId == uuid.Nil {
		err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!

			errTimes := 0
			if len(account.LoginErrorLock) > 0 {
				strErrTimes := string(account.LoginErrorLock[0])

				if strErrTimes != "" {
					errTimes, _ = strconv.Atoi(strErrTimes)
				}
			}

			errTimes = errTimes + 1
			strErrTimes := strconv.Itoa(errTimes)
			lastTime := strconv.FormatInt(time.Now().Unix(), 10)

			newStrLoginErrorLock := strErrTimes + "x" + lastTime
			loginLockUpdate := LoginLockUpdate(account.UserId, newStrLoginErrorLock)
			if !loginLockUpdate["status"].(bool) {
				return loginLockUpdate
			}

			return u.Message(false, "0x11022:Invalid login credentials. Please try again")
		}

		if len(account.LoginErrorLock) > 0 {

			strErrLockTimes := string(account.LoginErrorLock[0])
			errLockTimes := 0
			if strErrLockTimes != "" {
				errLockTimes, _ = strconv.Atoi(strErrLockTimes)
			}

			if errLockTimes < 9 {
				loginLockUpdate := LoginLockUpdate(account.UserId, "")
				if !loginLockUpdate["status"].(bool) {
					return loginLockUpdate
				}
			}

		}
	} else {
		if tokenUserId != account.UserId {
			LockAccount := LockAccount(account.UserId)
			if !LockAccount["status"].(bool) {
				return LockAccount
			}
			return u.Message(false, "0x11141:Rather than a technical error, a request was sent to the server with manually manipulated data. This is a fraud operation")
		}
	}

	roleId := u.None

	if email == os.Getenv("ROOT_ACCOUNT") {
		account.Nickname = "root"
		account.CreatedAt = time.Date(2018, time.September, 28, 0, 0, 0, 0, time.Local)
		account.AboutMe = "I have root rights"
		subject = "all rights reserved"
		roleId = u.Root
		// TOKEN_ROOT_SECRET_KEY
	}

	//Worked! 0x11023:Logged In
	account.Password = ""
	account.ID = 0

	atClaims := jwt.StandardClaims{}
	atClaims.Issuer = os.Getenv("JWT_ISSUER")
	atClaims.ExpiresAt = time.Now().Add(time.Minute * 1).Unix()
	atClaims.Subject = subject

	//Create JWT token
	tk := &u.Token{UserId: account.UserId, RoleId: roleId, StandardClaims: atClaims}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS512"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_SECRET_KEY")))

	if email == os.Getenv("ROOT_ACCOUNT") {
		tokenString, _ = token.SignedString([]byte(os.Getenv("TOKEN_ROOT_SECRET_KEY")))
	}

	account.Token = tokenString //Store the token in the response

	resp := u.Message(true, "0x11023:Logged In")
	resp["account"] = account
	return resp
}

func CheckUser(userId uuid.UUID, roleId int, appId int, merchantId uuid.UUID) map[string]interface{} {
	result := &ResultBuildToken{}

	err := d.GetDB().Table("membership.roles").Where("active and deleted_at is null and user_id = ? and app_id = ? and role_id = ? and merchant_id = ?", userId, appId, roleId, merchantId).Find(result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "0x11024:The user's role was not found for the specified application")
		}
		return u.Message(false, "0x11004:Connection error. Please retry")
	}

	if result.UserId != userId {
		return u.Message(false, "0x11025:The user's role was not match for the specified application")
	}

	resp := u.Message(true, "0x11026:Check user is ok")
	resp["account"] = result
	return resp
}

func BuildToken(userId uuid.UUID, roleId int, appId int, merchantId uuid.UUID, hasId bool, customData string) map[string]interface{} {

	resultCheck := CheckUser(userId, roleId, appId, merchantId)

	if !u.CheckOk(resultCheck) {
		return resultCheck
	}

	result := resultCheck["account"].(*ResultBuildToken)

	projectId := u.GetProjectId(result.ProjectInfo)

	// if errProjectId != nil {
	// 	return u.Message(false, "project info error")
	// }   goutils de bu olmali mi

	var initCompleted = !result.LastLoginDate.IsZero()

	account := &BuildAccount{}

	atClaims := jwt.StandardClaims{}
	atClaims.Issuer = os.Getenv("JWT_ISSUER")
	atClaims.ExpiresAt = time.Now().Add(time.Hour * 100).Unix()

	//Create JWT token
	tk := &u.Token{UserId: userId, RoleId: roleId, AppId: appId, MerchantId: merchantId, HasId: hasId, ProjectId: projectId, CustomData: customData, InitCompleted: initCompleted, StandardClaims: atClaims}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS512"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_SECRET_KEY")))
	account.Token = tokenString //Store the token in the response
	account.UserId = userId
	account.LastLoginDate = result.LastLoginDate

	if !hasId {
		err := d.GetDB().Exec("UPDATE membership.roles set last_login_date = ? where user_id = ? and app_id = ? and merchant_id = ?",
			time.Now(), result.UserId, appId, merchantId).Error
		if err != nil {
			return u.Message(false, "0x11027:Login success but, Connection error. Please retry")
		}
	}

	resp := u.Message(true, "0x11023:Logged In")
	resp["account"] = account
	return resp
}

func LoginLockUpdate(userId uuid.UUID, strLoginErrorLock string) map[string]interface{} {
	err := d.GetDB().Exec("UPDATE security.accounts set login_error_lock = ? where user_id = ?",
		strLoginErrorLock, userId).Error
	if err != nil {
		return u.Message(false, "0x11142:Connection error. Please retry. The update for account locking and incorrect login attempts could not be completed.")
	}
	return u.Message(true, "0x11031:Success")
}

func LockAccount(userId uuid.UUID) map[string]interface{} {
	err := d.GetDB().Exec("UPDATE security.accounts set lock = ? where user_id = ?",
		true, userId).Error
	if err != nil {
		return u.Message(false, "0x11142:Connection error. Please retry. The update for account locking and incorrect login attempts could not be completed.")
	}

	return u.Message(true, "0x11031:Success")
}
