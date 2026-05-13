// Copyright (c) 2017-2026 Onur Yaşar
// Licensed under AGPL v3 + Commercial Exception
// See LICENSE.txt

// https://github.com/rymory/rymory-core
// rymory.org 
// onuryasar.org
// onxorg@proton.me 

package zombie

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	u "github.com/rymory/goutils/api"
	d "github.com/rymory/goutils/db"
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

func GetZombieRoles() (map[string]interface{}, bool) {

	acc := &[]Role{}
	err := d.GetDB().Table("membership.roles").Joins("inner join application.projects on roles.merchant_id = projects.merchant_id").Where("projects.app_id = roles.app_id and ((application.projects.active = FALSE and application.projects.deleted_at is null and membership.roles.deleted_at is null) or (application.projects.deleted_at is not null))").Scan(acc).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "0x11004:Connection error. Please retry"), false
	}

	response := u.Message(true, "0x11048:Zombie roles found it")
	response["roles"] = acc
	return response, true
}
