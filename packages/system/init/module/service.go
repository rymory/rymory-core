// Copyright (c) 2017-2026 Onur Yaşar
// Licensed under AGPL v3 + Commercial Exception
// See LICENSE.txt

package initialize

import (
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	u "github.com/lemoras/goutils/api"
	d "github.com/lemoras/goutils/db"
)

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

type Role struct {
	gorm.Model
	//ID             int64     `gorm:"primaryKey"`
	UserId        uuid.UUID `json:"userId"`
	RoleId        int       `json:"roleId"`
	AppId         int       `json:"appId"`
	MerchantId    uuid.UUID `json:"merchantId"`
	Active        bool      `json:"active"`
	ProjectInfo   string    `json:"projectInfo"`
	LastLoginDate time.Time `json:"lastLoginDate"`
	CreatedBy     uuid.UUID `json:"createdBy"`
	UpdatedBy     uuid.UUID `json:"updatedBy"`
	DeletedBy     uuid.UUID `json:"deletedBy"`
}

type Project struct {
	gorm.Model
	Domains        string    `json:"domains"`
	AppId          int       `json:"appId"`
	MerchantId     uuid.UUID `json:"merchantId"`
	Active         bool      `json:"active"`
	ProjectName    string    `json:"projectName"`
	ManageMember   bool      `json:"manageMember"`
	IsOwnerLemoras bool      `json:"isOwnerLemoras"`
	CreatedBy      uuid.UUID `json:"createdBy"`
	UpdatedBy      uuid.UUID `json:"updatedBy"`
	DeletedBy      uuid.UUID `json:"deletedBy"`
}

type ValidateRole struct {
	UserId  uuid.UUID
	Role    bool
	Project string
	Check   bool
}

// Validate incoming user details...
func (account *Account) Validate() (map[string]interface{}, bool) {

	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "0x11001:Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "0x11002:Password is required"), false
	}

	//Email must be unique
	temp := &Account{}

	//check for errors and duplicate emails
	err := d.GetDB().Table("security.accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return u.Message(false, "0x11005:Email address already in use by another user."), false
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

	generatedPassword := ""

	if account.Email == os.Getenv("ROOT_ACCOUNT") {

		account.Nickname = "root"
		account.CreatedAt = time.Date(2018, time.September, 28, 0, 0, 0, 0, time.Local)
		account.AboutMe = "I have root rights"
		generatedPassword = u.GeneratePassword(20, true, true, true)
		account.Password = generatedPassword
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)

	account.Password = string(hashedPassword)
	account.UserId = uuid.New()

	d.GetDB().Table("security.accounts").Create(account)

	if account.ID <= 0 {
		return u.Message(false, "0x11007:Failed to create account, connection error.")
	}

	account.Token = ""

	account.Password = generatedPassword
	account.ID = 0

	response := u.Message(true, "0x11008:Account has been created")
	response["account"] = account
	return response
}

// Validate incoming user details...
func (project *Project) Validate() (map[string]interface{}, bool) {

	//Email must be unique
	temp := &Project{}
	//check for errors and duplicate emails
	err := d.GetDB().Table("application.projects").Where("app_id=? and merchant_id=?", project.AppId, project.MerchantId).Find(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry"), false
	}

	if err != nil && err == gorm.ErrRecordNotFound {
		return u.Message(true, "0x11006:Requirement passed"), true
	}

	return u.Message(false, "The project already has a record."), false
}

func (project *Project) Create(userId uuid.UUID) map[string]interface{} {

	if resp, ok := project.Validate(); !ok {
		return resp
	}

	project.Active = true
	project.CreatedBy = userId

	d.GetDB().Table("application.projects").Create(project)

	if project.ID <= 0 {
		return u.Message(false, "Failed to create project, connection error.")
	}

	project.ID = 0

	response := u.Message(true, "Project has been created")
	response["project"] = project
	return response
}

// Validate incoming user details...
func (role *Role) Validate(email string, tokenRole *Role, requiredRoleId int) (map[string]interface{}, bool, uuid.UUID, string) {
	//Email must be unique
	temp := &ValidateRole{}
	//check for errors and duplicate emails
	err := d.GetDB().Raw("select security.accounts.user_id, EXISTS(select 1 from membership.roles where roles.user_id= ? and roles.app_id=? and roles.merchant_id = ? and roles.active and roles.deleted_at is null) as check , EXISTS(select 1 from membership.roles where roles.user_id= security.accounts.user_id and roles.app_id =? and roles.merchant_id = ?) as role, (select string_agg(CONCAT (application.projects .id::TEXT, ':', application.projects.project_name::TEXT), ', ') from application.projects where projects.manage_member = FALSE and projects.active and projects.deleted_at is null and projects.app_id =? and projects.merchant_id = ?) as project from security.accounts  where accounts.email = ?", tokenRole.UserId, tokenRole.AppId, tokenRole.MerchantId, role.AppId, role.MerchantId, role.AppId, role.MerchantId, email).Find(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry"), false, temp.UserId, temp.Project
	}

	if err != nil && err == gorm.ErrRecordNotFound {
		return u.Message(false, "Email Not found. Please retry"), false, temp.UserId, temp.Project
	}

	if !temp.Check && !(tokenRole.RoleId == u.Root && requiredRoleId == u.MerchantAdmin) {
		return u.Message(false, "The current user/token not expire but it hasnt right auth."), false, temp.UserId, temp.Project
	}

	if len(temp.Project) == 0 {
		return u.Message(false, "The App and Merchant not found in this role for this application or project not available."), false, temp.UserId, temp.Project
	}

	if temp.Role {
		return u.Message(false, "The user already has a record in this role for this application."), false, temp.UserId, temp.Project
	}
	return u.Message(true, "0x11006:Requirement passed"), true, temp.UserId, temp.Project
}

func (role *Role) Create(email string, tokenRole *Role, requiredRoleId int) map[string]interface{} {

	resp, ok, userId, project := role.Validate(email, tokenRole, requiredRoleId)
	if !ok {
		return resp
	}

	role.UserId = userId
	role.ProjectInfo = project

	d.GetDB().Table("membership.roles").Create(role)

	if role.ID <= 0 {
		return u.Message(false, "Failed to create role, connection error.")
	}

	role.ID = 0

	response := u.Message(true, "Role has been created")
	response["role"] = role
	return response
}

func GetIsOwnLemorasMerchant() (map[string]interface{}, uuid.UUID, bool) {

	//Email must be unique
	temp := &Project{}
	//check for errors and duplicate emails
	err := d.GetDB().Table("application.projects").Where("is_owner_lemoras").First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry"), temp.MerchantId, false
	}

	if err != nil && err == gorm.ErrRecordNotFound {
		return u.Message(false, "Project was not found by IsOwnLemoras"), temp.MerchantId, false
	}

	return u.Message(true, "The project found."), temp.MerchantId, true
}
