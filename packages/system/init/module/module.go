package initialize

import (
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	u "gitlab.com/onxorg/goutils/api"
	d "gitlab.com/onxorg/goutils/db"
)

type Request struct {
	Http     CustomHttp       `json:"http"`
	Projects []ProjectRequest `json:"projects"`
}

type CustomHttp struct {
	CustomHeader CustomHeader `json:"headers"`
	Method       string       `json:"method"`
	Path         string       `json:"path"`
}

type CustomHeader struct {
	Authorization string `json:"authorization"`
}

type ProjectRequest struct {
	AppId   int      `json:"appId"`
	Domains []string `json:"domains"`
	AppName string   `json:"appName"`
}

func init() {
	//MigrationModels()
}

func Invoke(in Request) (*u.Response, error) {

	resp := u.Message(false, "0x11028:Invalid request")

	path := strings.Replace(in.Http.Path, "/", "", -1)

	if path == "nemutluturkumdiyene" {
		MigrationModels()
		resp = CreateRootAccount(in.Projects)
	} else if path == "update" {
		context := &u.Context{}
		if isOk, res := u.JwtAuthentication(in.Http.CustomHeader.Authorization, context); !isOk {
			return &res, nil
		}
		if context.RoleId == u.Root {
			MigrationModels()
			resp = u.Message(true, "0x11034:Migration done..")
			if len(in.Projects) > 0 {
				rootEmail := strings.ToLower(os.Getenv("root_account"))
				createNewProject(rootEmail, context.UserId, in.Projects)
				resp = u.Message(true, "0x11035:Migration & new project done..")
			}
		}
	}
	return u.Respond(resp)
}

var CreateRootAccount = func(projects []ProjectRequest) map[string]interface{} {

	rootAccount := &Account{}

	rootAccount.Email = strings.ToLower(os.Getenv("root_account"))
	rootAccount.Password = "12345678"

	rootResp := rootAccount.Create()

	if !u.CheckOk(rootResp) {
		return rootResp
	}

	merchantId := uuid.New()
	//demoUserId := uuid.New()
	for _, projectItem := range projects {
		domains := ""
		if len(projectItem.Domains) > 0 {
			domains = strings.Join(projectItem.Domains, ",")
		}
		appName := strings.ToTitle(strings.ToLower(projectItem.AppName))
		project := &Project{}
		project.IsOwnerLemoras = true
		project.AppId = projectItem.AppId
		project.Domains = domains
		project.MerchantId = merchantId
		project.ProjectName = appName + " project"

		rootProject := project.Create(rootAccount.UserId)

		if !u.CheckOk(rootProject) {
			return rootProject
		}
	}

	account := &Account{}
	account.Email = strings.ToLower("madmin@lemoras.com")
	account.Password = "123456Wqm"

	mAdminAccount := account.Create()

	if !u.CheckOk(mAdminAccount) {
		return mAdminAccount
	}

	for _, roleItem := range projects {

		role := &Role{}

		appName := strings.ToTitle(strings.ToLower(roleItem.AppName))

		role.ProjectInfo = string(rune(role.ID)) + ":" + appName + " project"
		role.AppId = roleItem.AppId
		role.RoleId = 2
		role.MerchantId = merchantId
		role.Active = true
		role.CreatedBy = rootAccount.UserId

		tokenRole := &Role{
			AppId:  roleItem.AppId,
			RoleId: 1,
		}

		mAdminrole := role.Create(account.Email, tokenRole, 2)
		if !u.CheckOk(mAdminrole) {
			return mAdminrole
		}
	}

	return rootResp
}

var MigrationModels = func() {

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "security." + defaultTableName
	}

	d.GetDB().Exec("CREATE SCHEMA IF NOT EXISTS security")

	d.GetDB().Debug().AutoMigrate(&Account{})

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "membership." + defaultTableName
	}

	d.GetDB().Exec("CREATE SCHEMA IF NOT EXISTS membership")

	d.GetDB().Debug().AutoMigrate(&Role{})

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "application." + defaultTableName
	}

	d.GetDB().Exec("CREATE SCHEMA IF NOT EXISTS application")

	d.GetDB().Debug().AutoMigrate(&Project{})
}

var createNewProject = func(email string, userId uuid.UUID, projects []ProjectRequest) map[string]interface{} {

	resp, merchantId, ok := GetIsOwnLemorasMerchant()

	if !ok {
		return resp
	}

	for _, projectItem := range projects {
		domains := ""
		if len(projectItem.Domains) > 0 {
			domains = strings.Join(projectItem.Domains, ",")
		}
		appName := strings.ToTitle(strings.ToLower(projectItem.AppName))
		project := &Project{}
		project.IsOwnerLemoras = true
		project.AppId = projectItem.AppId
		project.Domains = domains
		project.MerchantId = merchantId
		project.ProjectName = appName + " project"

		project.Create(userId)
	}

	for _, roleItem := range projects {

		role := &Role{}

		appName := strings.ToTitle(strings.ToLower(roleItem.AppName))

		role.ProjectInfo = string(rune(role.ID)) + ":" + appName + " project"
		role.AppId = roleItem.AppId
		role.RoleId = 2
		role.MerchantId = merchantId
		role.Active = true
		role.CreatedBy = userId

		tokenRole := &Role{
			AppId:  roleItem.AppId,
			RoleId: 1,
		}

		role.Create(email, tokenRole, 2)
	}

	return u.Message(true, "0x11036:Created new projects")
}
