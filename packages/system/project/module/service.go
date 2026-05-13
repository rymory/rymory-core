// Copyright (c) 2017-2026 Onur Yaşar
// Licensed under AGPL v3 + Commercial Exception
// See LICENSE.txt

// https://github.com/rymory/rymory-core
// rymory.org 
// onuryasar.org
// onxorg@proton.me 

package project

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	u "github.com/rymory/goutils/api"
	d "github.com/rymory/goutils/db"
)

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

type RequiredProject struct {
	AppId          int       `json:"appId"`
	MerchantId     uuid.UUID `json:"merchantId"`
	RequiredStatus int       `json:"requiredStatus"`
	ProjectName    string    `json:"projectName"`
	ManageMember   int       `json:"manageMember"`
	MemberStatus   int       `json:"memberStatus"`
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

type GetRoleCheck struct {
	Check bool `json:"check"`
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

func GetProject(app int, merchant uuid.UUID) (*Project, bool) {

	acc := &Project{}
	err := d.GetDB().Table("application.projects").Where("app_id =? and merchant_id = ?", app, merchant).First(acc).Error
	if err != nil {
		return nil, false
	}
	return acc, true
}

func (requiredProject *RequiredProject) GetProject() (*Project, bool) {

	return GetProject(requiredProject.AppId, requiredProject.MerchantId)
}

func GetProjectsByMerchantId(merchantId uuid.UUID) (map[string]interface{}, bool) {

	acc := &[]Project{}
	var err error
	if merchantId == uuid.Nil {
		err = d.GetDB().Table("application.projects").Scan(acc).Error
	} else {
		err = d.GetDB().Table("application.projects").Where("deleted_at is null and merchant_id=?", merchantId).Scan(acc).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry"), false
	}
	response := u.Message(true, "Projects found it")
	response["projects"] = acc
	return response, true
}

func GetProjectByDomain(domain string) (map[string]interface{}, bool) {

	acc := &[]Project{}
	domainParam := "%" + domain + "%"
	err := d.GetDB().Table("application.projects").Where("deleted_at is null and active and domains like ?", domainParam).Scan(acc).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry"), false
	}

	if err != nil && err == gorm.ErrRecordNotFound {
		return u.Message(false, "Project not found"), false
	}
	jsonData, _ := json.Marshal(acc)

	data := []Project{}
	json.Unmarshal(jsonData, &data)
	foundState := false
	resData := []map[string]interface{}{}
	for _, p := range data {
		var tmpDomains = strings.Split(p.Domains, ",")
		for _, item := range tmpDomains {
			if item == domain {
				check := true
				for _, v := range resData {
					if v["projectId"] == p.ID {
						check = false
					}
				}
				if check {
					foundState = true
					tmpData := map[string]interface{}{"merchantId": p.MerchantId, "manageMember": p.ManageMember, "projectId": p.ID, "appId": p.AppId}
					resData = append(resData, tmpData)
				}
			}
		}
	}

	if foundState {
		response := u.Message(true, "Projects found it")
		response["projects"] = resData
		return response, true
	}

	return u.Message(false, "Project not found"), false
}

func (requiredProject *RequiredProject) Update(user uuid.UUID) (map[string]interface{}, *Project, bool) {

	project, ok := GetProject(requiredProject.AppId, requiredProject.MerchantId)

	if !ok {
		return u.Message(false, "0x11039:Project not found or already soft deleted"), project, false
	}

	if requiredProject.ProjectName != "" {
		project.ProjectName = requiredProject.ProjectName
	}

	if requiredProject.ManageMember == 1 {
		project.ManageMember = false
	}

	if requiredProject.ManageMember == 2 {
		project.ManageMember = true
	}

	if requiredProject.RequiredStatus == 1 {
		project.Active = false
	}

	if requiredProject.RequiredStatus == 2 {
		project.Active = true
	}

	//project.Active = !project.Active
	project.UpdatedBy = user

	err := d.GetDB().Table("application.projects").Where("deleted_at is null").Save(project).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry"), project, false
	}
	var msg string
	if project.Active {
		msg = "activated"
	} else {
		msg = "passived"
	}
	return u.Message(true, "Project has been "+msg), project, true
}

func (project *Project) Delete(user uuid.UUID) (map[string]interface{}, bool) {

	//project = GetProject(project.AppId, project.MerchantId)

	project.DeletedBy = user
	err := d.GetDB().Table("application.projects").Model(&project).Update(&project).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry"), false
	}

	err = d.GetDB().Table("application.projects").Delete(project).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry"), false
	}

	return u.Message(true, "Project has been deleted"), true
}

func (project *Project) RevertDelete() (map[string]interface{}, bool) {

	err := d.GetDB().Exec("UPDATE FROM application.projects set deleted_at =?, deleted_by=? WHERE appId =? and merchantId=?", nil, nil, project.AppId, project.MerchantId).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry, revert delete project"), false
	}

	return u.Message(true, "0x11006:Requirement passed"), true
}

func BulkRoleDelete(appId int, merchantId uuid.UUID, isSoftDelete bool) (map[string]interface{}, bool) {
	var err error

	if isSoftDelete {
		err = d.GetDB().Table("membership.roles").Where("app_id =? and merchant_id=?", appId, merchantId).Delete(&Role{}).Error
	} else {
		err = d.GetDB().Table("membership.roles").Where("app_id =? and merchant_id=?", appId, merchantId).Unscoped().Delete(&Role{}).Error // meybe history table record
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry"), false
	}

	return u.Message(true, "0x11006:Requirement passed"), true
}

func BulkRoleRevertDelete(appId int, merchantId uuid.UUID) (map[string]interface{}, bool) {

	err := d.GetDB().Exec("UPDATE membership.roles set deleted_at = null, deleted_by = null WHERE app_id =? and merchant_id=?", appId, merchantId).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry, revert delete project"), false
	}

	return u.Message(true, "0x11006:Requirement passed"), true
}

func BulkRoleMember(appId int, merchantId uuid.UUID, status int) (map[string]interface{}, bool) {

	var err error
	active := false

	if status == 2 {
		active = true
	}

	if status == 3 {
		err = d.GetDB().Exec("DELETE from membership.roles WHERE role_id = ? and app_id =? and merchant_id=?", u.Member, appId, merchantId).Error
	} else {
		err = d.GetDB().Exec("UPDATE membership.roles set active = ? WHERE role_id = ? and app_id =? and merchant_id=?", active, u.Member, appId, merchantId).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry, revert delete project"), false
	}

	return u.Message(true, "0x11006:Requirement passed"), true
}

func BulkRoleUpdateProjectInfo(projectInfo string, appId int, merchantId uuid.UUID) (map[string]interface{}, bool) {

	err := d.GetDB().Exec("UPDATE membership.roles set project_info= ? WHERE app_id =? and merchant_id=?", projectInfo, appId, merchantId).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry, revert delete project"), false
	}

	return u.Message(true, "0x11006:Requirement passed"), true
}

func GetRole(userId uuid.UUID, app int, merchant uuid.UUID, roleId int) (map[string]interface{}, bool) {

	acc := &GetRoleCheck{}
	err := d.GetDB().Raw("select EXISTS(select 1 from membership.roles inner join application.projects on roles.merchant_id = projects.merchant_id where projects.active and projects.deleted_at is null and projects.app_id = roles.app_id and roles.user_id = ? and roles.app_id =? and roles.merchant_id = ? and roles.role_id = ? and roles.active and roles.deleted_at is null) as check", userId, app, merchant, roleId).First(acc).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "0x11024:The user's role was not found for the specified application"), false
		}
		return u.Message(false, "0x11004:Connection error. Please retry"), false
	}

	if !acc.Check && !(roleId == u.Root) {
		return u.Message(false, "The current user/token not expire but it hasnt right auth."), false
	}

	return u.Message(true, "0x11006:Requirement passed"), true
}
