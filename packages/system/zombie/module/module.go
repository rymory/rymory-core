// Copyright (c) 2017-2026 Onur Yaşar
// Licensed under AGPL v3 + Commercial Exception
// See LICENSE.txt

// https://github.com/rymory/rymory-core
// rymory.org 
// onuryasar.org
// onxorg@proton.me 

package zombie

import (
	u "github.com/rymory/goutils/api"
)

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	PhotoUrl string `json:"photoUrl"`
	AboutMe  string `json:"aboutMe"`

	Http u.CustomHttp `json:"http"`
}

func Invoke(in Request) (*u.Response, error) {

	context := &u.Context{}
	if isOk, res := u.JwtAuthentication(in.Http.CustomHeader.Authorization, context); !isOk {
		return &res, nil
	}

	return u.Respond(GetAllZombieRole(*context))
}

var GetAllZombieRole = func(context u.Context) map[string]interface{} {

	roleId := context.RoleId
	tokenRoleLevel := u.GetRoleLevel(roleId)
	if tokenRoleLevel == u.Root {
		resp, _ := GetZombieRoles()
		return resp
	}

	return u.Message(false, "0x11020:You do not have access authority")
}
