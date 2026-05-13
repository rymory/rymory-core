// Copyright (c) 2017-2026 Onur Yaşar
// Licensed under AGPL v3 + Commercial Exception
// See LICENSE.txt

package member

import (
	"github.com/google/uuid"
	u "github.com/lemoras/goutils/api"
)

type Request struct {
	AppId      int       `json:"appId"`
	MerchantId uuid.UUID `json:"merchantId"`

	Http u.CustomHttp `json:"http"`
}

func Invoke(in Request) (*u.Response, error) {

	var resp map[string]interface{}

	context := &u.Context{}
	if isOk, res := u.JwtAuthentication(in.Http.CustomHeader.Authorization, context); !isOk {
		return &res, nil
	}

	switch in.Http.Method {
	case "GET":
		resp = GetAllMember(in, *context)
		break
	case "POST":
		resp = PostAllMember(in, *context)
		break
	default:
		resp = u.Message(false, "0x11028:Invalid request")
		break
	}

	return u.Respond(resp)
}

var GetAllMember = func(request Request, context u.Context) map[string]interface{} {

	request.AppId = context.AppId
	request.MerchantId = context.MerchantId

	return PostAllMember(request, context)
}

var PostAllMember = func(request Request, context u.Context) map[string]interface{} {

	appId := context.AppId
	merchantId := context.MerchantId
	roleId := context.RoleId

	tokenRoleLevel := u.GetRoleLevel(roleId)

	if tokenRoleLevel == u.Root {

		if request.AppId == 0 {
			return u.Message(false, "0x11018:AppId is zero value")
		}

		if request.MerchantId == uuid.Nil {
			return u.Message(false, "0x11019:MerchantId is zero value")
		}
	} else if tokenRoleLevel == u.MerchantAdmin {
		if request.AppId == 0 {
			request.AppId = appId
			//return u.Message(false, "0x11018:AppId is zero value")
		}
		request.MerchantId = merchantId
	} else if tokenRoleLevel > u.MerchantAdmin && tokenRoleLevel < u.User {
		request.AppId = appId
		request.MerchantId = merchantId
	} else {
		return u.Message(false, "0x11020:You do not have access authority")
	}

	respFail, ok, members := GetMembers(request.AppId, request.MerchantId, roleId)
	if !ok {
		return respFail
	}

	resp := make(map[string]interface{})
	resp["members"] = members
	return resp
}
