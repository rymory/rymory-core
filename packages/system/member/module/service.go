// Copyright (c) 2017-2026 Onur Yaşar
// Licensed under AGPL v3 + Commercial Exception
// See LICENSE.txt

// https://github.com/rymory/rymory-core
// rymory.org 
// onuryasar.org
// onxorg@proton.me 

package member

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	u "github.com/rymory/goutils/api"
	d "github.com/rymory/goutils/db"
)

type Member struct {
	UserId        uuid.UUID `json:"userId"`
	RoleId        int       `json:"roleId"`
	AppId         int       `json:"appId"`
	MerchantId    uuid.UUID `json:"merchantId"`
	Active        bool      `json:"active"`
	ProjectInfo   string    `json:"projectInfo"`
	ProjectName   string    `json:"projectName"`
	LastLoginDate time.Time `json:"lastLoginDate"`
	CreatedByName string    `json:"createdByName"`
	UpdatedByName string    `json:"updatedByName"`
	DeletedByName string    `json:"deletedByName"`
	Nickname      string    `json:"nickname"`
	PhotoUrl      string    `json:"photoUrl"`
	AboutMe       string    `json:"aboutMe"`
	Email         string    `json:"email"`
	UpdatedAt     time.Time `json:"updatedAt"`
	CreatedAt     time.Time `json:"createdAt"`
}

func GetMembers(app int, merchant uuid.UUID, role int) (map[string]interface{}, bool, *[]Member) {
	role = u.GetRoleLevel(role)

	acc := &[]Member{}
	err := d.GetDB().Raw(`select 
			roles.user_id,
			roles.role_id,
			roles.app_id,
			roles.merchant_id,
			roles.active,
			roles.project_info,
			projects.project_name,
			roles.last_login_date,
			(select nickname from security.accounts where user_id = roles.created_by) as created_by_name,
			(select nickname from security.accounts where user_id = roles.updated_by) as updated_by_name,
			(select nickname from security.accounts where user_id = roles.deleted_by) as deleted_by_name,
			accounts.nickname,
			accounts.photo_url,
			accounts.about_me,
			accounts.email,
			roles.created_at,
			roles.updated_at
		 from membership.roles 
		 inner join application.projects 
		 	on roles.merchant_id = projects.merchant_id 
		 inner join security.accounts 
		 	on roles.user_id = accounts.user_id 
		where projects.active and projects.deleted_at is null 
			and projects.app_id = roles.app_id
			and roles.app_id =? and roles.merchant_id = ?
			and cast(cast(roles.role_id as varchar(1)) as smallint) > ? and cast(cast(roles.role_id as varchar(1)) as smallint) < ? `,
		app, merchant, role, u.Member).Scan(acc).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "0x11024:The user's role was not found for the specified application"), false, acc
		}
		return u.Message(false, "0x11004:Connection error. Please retry"), false, acc
	}

	return u.Message(true, "0x11006:Requirement passed"), true, acc
}
