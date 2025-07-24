package role

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	u "gitlab.com/onxorg/goutils/api"
	d "gitlab.com/onxorg/goutils/db"
)

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

type GetRoleCheck struct {
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
	Check         bool      `json:"check"`
}

type ValidateRoleModel struct {
	UserId  uuid.UUID
	Role    bool
	Project string
	Check   bool
}

type SelfValidateRoleModel struct {
	UserId  uuid.UUID
	Role    bool
	Project string
}

// Validate incoming user details...
func (role *Role) Validate(email string, tokenRole *Role, requiredRoleId int) (map[string]interface{}, bool, uuid.UUID, string) {
	//Email must be unique
	temp := &ValidateRoleModel{}
	//check for errors and duplicate emails and without lock (no lock account) where condition and account.lock = ? (false)
	err := d.GetDB().Raw("select security.accounts.user_id, EXISTS(select 1 from membership.roles where roles.role_id != ? and roles.user_id= ? and roles.app_id=? and roles.merchant_id = ? and roles.active and roles.deleted_at is null) as check , EXISTS(select 1 from membership.roles where roles.user_id= security.accounts.user_id and roles.app_id =? and roles.merchant_id = ?) as role, (select string_agg(CONCAT (application.projects.id::TEXT, ':', application.projects.project_name::TEXT, ':', application.projects.manage_member::INTEGER)), ', ') from application.projects where projects.active and projects.deleted_at is null and projects.app_id =? and projects.merchant_id = ?) as project from security.accounts  where accounts.email = ? and accounts.lock = ?", u.Member, tokenRole.UserId, tokenRole.AppId, tokenRole.MerchantId, role.AppId, role.MerchantId, role.AppId, role.MerchantId, email, false).Find(temp).Error
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

	tmpProjectInfos := strings.Split(temp.Project, ":")
	temp.Project = tmpProjectInfos[0] + ":" + tmpProjectInfos[1]

	if requiredRoleId == u.Member && tmpProjectInfos[2] != "1" {
		return u.Message(false, "The App and Merchant not found in this role for this application or project not available."), false, temp.UserId, temp.Project
	}

	if temp.Role {
		return u.Message(false, "The user already has a record in this role for this application."), false, temp.UserId, temp.Project
	}
	return u.Message(true, "0x11006:Requirement passed"), true, temp.UserId, temp.Project
}

func (role *Role) SelfValidate(userId uuid.UUID) (map[string]interface{}, bool, string) {

	//Email must be unique
	temp := &SelfValidateRoleModel{}
	//check for errors and duplicate emails and without lock (no lock account) where condition and account.lock = ? (false)
	err := d.GetDB().Raw("select security.accounts.user_id, EXISTS(select 1 from membership.roles where roles.role_id = ? and roles.user_id= security.accounts.user_id and roles.app_id =? and roles.merchant_id = ?) as role, (select string_agg(CONCAT (application.projects .id::TEXT, ':', application.projects.project_name::TEXT), ', ') from application.projects where projects.manage_member and projects.active and projects.deleted_at is null and projects.app_id =? and projects.merchant_id = ?) as project  from security.accounts where accounts.user_id = ? and accounts.lock = ?", u.Member, role.AppId, role.MerchantId, role.AppId, role.MerchantId, userId, false).Find(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry"), false, temp.Project
	}

	if err != nil && err == gorm.ErrRecordNotFound {
		return u.Message(false, "User Not found. Please retry"), false, temp.Project
	}

	if len(temp.Project) == 0 {
		return u.Message(false, "The App and Merchant not found in this role for this application or project not available or not accepted member role."), false, temp.Project
	}

	if temp.Role {
		return u.Message(false, "The user already has a record in this role for this application."), false, temp.Project
	}

	return u.Message(true, "0x11006:Requirement passed"), true, temp.Project
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

func (role *Role) SelfCreate(userId uuid.UUID) map[string]interface{} {

	resp, ok, project := role.SelfValidate(userId)
	if !ok {
		return resp
	}

	role.ProjectInfo = project
	role.UserId = userId
	role.CreatedBy = userId

	d.GetDB().Table("membership.roles").Create(role)

	if role.ID <= 0 {
		return u.Message(false, "Failed to create role, connection error.")
	}

	role.ID = 0

	response := u.Message(true, "Role has been created")
	response["role"] = role
	return response
}

func GetRole(userId uuid.UUID, app int, merchant uuid.UUID, roleId int, tokenRole *Role) (map[string]interface{}, bool, *Role) {

	acc := &GetRoleCheck{}
	res := &Role{}
	err := d.GetDB().Raw("select EXISTS(select 1 from membership.roles where roles.role_id = ? and roles.user_id = ? and roles.app_id= ? and roles.merchant_id = ? and roles.active and roles.deleted_at is null) as check, * from membership.roles inner join application.projects on roles.merchant_id = projects.merchant_id where projects.active and projects.deleted_at is null and projects.app_id = roles.app_id and roles.user_id = ? and roles.app_id =? and roles.merchant_id = ? and roles.role_id = ?", tokenRole.RoleId, tokenRole.UserId, tokenRole.AppId, tokenRole.MerchantId, userId, app, merchant, roleId).First(acc).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "0x11024:The user's role was not found for the specified application"), false, res
		}
		return u.Message(false, "0x11004:Connection error. Please retry"), false, res
	}

	if !acc.Check && !(tokenRole.RoleId == u.Root && acc.RoleId == u.MerchantAdmin) {
		return u.Message(false, "The current user/token not expire but it hasnt right auth."), false, res
	}

	res.UserId = acc.UserId
	res.RoleId = acc.RoleId
	res.AppId = acc.AppId
	res.MerchantId = acc.MerchantId
	res.Active = acc.Active
	res.ProjectInfo = acc.ProjectInfo
	res.LastLoginDate = acc.LastLoginDate
	res.CreatedBy = acc.CreatedBy
	res.UpdatedBy = acc.UpdatedBy

	return u.Message(true, "0x11006:Requirement passed"), true, res
}

// asagidakimetodu sil yukarisini kullan
func (updateRole *Request) GetRole(tokenRole *Role) (map[string]interface{}, bool, *Role) {

	return GetRole(updateRole.UserId, updateRole.AppId, updateRole.MerchantId, updateRole.RoleId, tokenRole)
}

func GetRolesById(uid uuid.UUID) (map[string]interface{}, bool) {

	acc := &[]Role{}
	err := d.GetDB().Table("membership.roles").Where("deleted_at is null and user_id = ?", uid).Scan(acc).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry"), false
	}
	response := u.Message(true, "Roles found it")
	response["roles"] = acc
	return response, true
}

func (role *Role) Update(user uuid.UUID) map[string]interface{} {

	// resp, ok, role := GetRole(role.UserId, role.AppId, role.MerchantId)
	// if !ok {
	// 	return resp
	// }
	role.Active = !role.Active
	role.UpdatedBy = user

	err := d.GetDB().Table("membership.roles").Unscoped().Model(&Role{}).Where("user_id=? and app_id=? and merchant_id=? and role_id = ?", role.UserId, role.AppId, role.MerchantId, role.RoleId).Updates(map[string]interface{}{"active": role.Active, "updated_by": role.UpdatedBy}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry")
	}
	var msg string
	if role.Active {
		msg = "activated"
	} else {
		msg = "passived"
	}
	response := u.Message(true, "Role has been "+msg)
	response["role"] = role
	return response
}

func (role *Role) Delete(user uuid.UUID) map[string]interface{} {

	// resp, ok, role := GetRole(role.UserId, role.AppId, role.MerchantId)
	// if !ok {
	// 	return resp
	// }

	role.DeletedBy = user
	err := d.GetDB().Table("membership.roles").Unscoped().Model(&Role{}).Where("user_id=? and app_id=? and merchant_id=?", role.UserId, role.AppId, role.MerchantId).Delete(&Role{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry")
	}

	response := u.Message(true, "Role has been deleted")
	response["role"] = role
	return response
}
