package role

import (
	"strings"

	"github.com/google/uuid"
	u "github.com/lemoras/goutils/api"
)

type Request struct {
	UserId     uuid.UUID `json:"userId"`
	Email      string    `json:"email"`
	RoleId     int       `json:"roleId"`
	AppId      int       `json:"appId"`
	MerchantId uuid.UUID `json:"merchantId"`
	Active     bool      `json:"active"`

	Http CustomHttp `json:"http"`
}

type CustomHttp struct {
	CustomHeader CustomHeader `json:"headers"`
	Method       string       `json:"method"`
}

type CustomHeader struct {
	Authorization string `json:"authorization"`
}

func Invoke(in Request) (*u.Response, error) {

	var resp map[string]interface{}

	context := &u.Context{}
	if isOk, res := u.JwtAuthentication(in.Http.CustomHeader.Authorization, context); !isOk {
		return &res, nil
	}

	switch in.Http.Method {
	case "GET":
		resp = GetRoles(*context)
		break
	case "POST":
		resp = CreateRole(in, *context)
		break
	case "PUT":
		resp = UpdateRole(in, *context)
		break
	case "DELETE":
		resp = DeleteRole(in, *context)
		break
	default:
		resp = u.Message(false, "0x11028:Invalid request")
		break
	}

	return u.Respond(resp)
}

var CreateRole = func(createRole Request, context u.Context) map[string]interface{} {

	if createRole.AppId == 0 {
		return u.Message(false, "0x11018:AppId is zero value")
	}

	if createRole.MerchantId == uuid.Nil {
		return u.Message(false, "0x11019:MerchantId is zero value")
	}

	userId := context.UserId
	roleId := context.RoleId

	appId := context.AppId
	merchantId := context.MerchantId

	tokenRoleLevel := u.GetRoleLevel(roleId)
	if appId == 0 && merchantId == uuid.Nil && tokenRoleLevel == u.None {
		userId := context.UserId
		role := &Role{}
		role.RoleId = u.Member
		role.Active = true
		role.AppId = createRole.AppId
		role.MerchantId = createRole.MerchantId
		return role.SelfCreate(userId) //Create account
	}

	if !strings.Contains(createRole.Email, "@") {
		return u.Message(false, "0x11001:Email address is required")
	}

	if createRole.RoleId == 0 {
		return u.Message(false, "0x11029:RoleId is zero value")
	}

	role := &Role{}
	role.RoleId = createRole.RoleId

	requiredRoleLevel := u.GetRoleLevel(role.RoleId)

	if (tokenRoleLevel == u.Superuser && requiredRoleLevel == u.User) || (tokenRoleLevel == u.Admin && requiredRoleLevel > u.Admin && requiredRoleLevel < u.Member) || (tokenRoleLevel == u.MerchantAdmin && requiredRoleLevel > u.MerchantAdmin && requiredRoleLevel < u.Member) || (tokenRoleLevel == u.Root && requiredRoleLevel > u.Root && requiredRoleLevel < u.Member) {

		role.AppId = appId
		role.MerchantId = merchantId
		role.Active = true
		role.CreatedBy = userId

		tokenRole := &Role{
			UserId:     userId,
			AppId:      appId,
			MerchantId: merchantId,
			RoleId:     tokenRoleLevel,
		}

		if tokenRoleLevel == u.Root {
			role.Active = false
			role.AppId = createRole.AppId
			role.MerchantId = createRole.MerchantId
		}

		createRole.Email = strings.ToLower(createRole.Email)
		return role.Create(createRole.Email, tokenRole, createRole.RoleId) //Create account

	} else {
		return u.Message(false, "0x11020:You do not have access authority")
	}
}

var UpdateRole = func(updateRole Request, context u.Context) map[string]interface{} {

	if updateRole.AppId == 0 {
		return u.Message(false, "0x11018:AppId is zero value")
	}

	if updateRole.MerchantId == uuid.Nil {
		return u.Message(false, "0x11019:MerchantId is zero value")
	}

	if updateRole.UserId == uuid.Nil {
		return u.Message(false, "0x11030:UserId is zero value")
	}

	if updateRole.RoleId == 0 {
		return u.Message(false, "0x11029:RoleId is zero value")
	}

	appId := context.AppId
	merchantId := context.MerchantId
	userId := context.UserId

	roleId := context.RoleId
	tokenRoleLevel := u.GetRoleLevel(roleId)

	if (appId == updateRole.AppId && merchantId == updateRole.MerchantId) || (tokenRoleLevel == u.MerchantAdmin && merchantId == updateRole.MerchantId) || (tokenRoleLevel == u.Root) {
		tokenRole := &Role{
			UserId:     userId,
			AppId:      appId,
			MerchantId: merchantId,
			RoleId:     roleId,
		}

		resp, ok, role := updateRole.GetRole(tokenRole)
		if !ok {
			return resp
		}

		requiredRoleLevel := u.GetRoleLevel(role.RoleId)

		if (tokenRoleLevel == u.Superuser && requiredRoleLevel > u.Superuser && requiredRoleLevel <= u.Member) || (tokenRoleLevel == u.Admin && requiredRoleLevel > u.Admin && requiredRoleLevel <= u.Member) || (tokenRoleLevel == u.MerchantAdmin && requiredRoleLevel > u.MerchantAdmin && requiredRoleLevel <= u.Member) || (tokenRoleLevel == u.Root && requiredRoleLevel > u.Root && requiredRoleLevel <= u.Member) {

			return role.Update(userId) //Update account
		} else {
			return u.Message(false, "0x11020:You do not have access authority")
		}
	} else {
		return u.Message(false, "0x11020:You do not have access authority")
	}
}

var GetRoles = func(context u.Context) map[string]interface{} {

	id := context.UserId
	resp, _ := GetRolesById(id) //Get account
	return resp
}

var DeleteRole = func(updateRole Request, context u.Context) map[string]interface{} {

	appId := context.AppId
	merchantId := context.MerchantId
	userId := context.UserId

	roleId := context.RoleId

	tokenRoleLevel := u.GetRoleLevel(roleId)

	if tokenRoleLevel == u.Member {
		updateRole.UserId = userId
		updateRole.AppId = appId
		updateRole.MerchantId = merchantId
		updateRole.RoleId = roleId

		tokenRole := &Role{
			UserId:     userId,
			AppId:      appId,
			MerchantId: merchantId,
			RoleId:     roleId,
		}
		respFail, ok, role := updateRole.GetRole(tokenRole)
		if !ok {
			return respFail
		}

		return role.Delete(userId) //Update account
	}

	if updateRole.AppId == 0 {
		return u.Message(false, "0x11018:AppId is zero value")
	}

	if updateRole.MerchantId == uuid.Nil {
		return u.Message(false, "0x11019:MerchantId is zero value")
	}

	if updateRole.UserId == uuid.Nil {
		return u.Message(false, "0x11030:UserId is zero value")
	}

	if updateRole.RoleId == 0 {
		return u.Message(false, "0x11029:RoleId is zero value")
	}

	if (appId == updateRole.AppId && merchantId == updateRole.MerchantId) || (tokenRoleLevel == u.MerchantAdmin && merchantId == updateRole.MerchantId) || (tokenRoleLevel == u.Root) {

		tokenRole := &Role{
			UserId:     userId,
			AppId:      appId,
			MerchantId: merchantId,
			RoleId:     roleId,
		}
		respFail, ok, role := updateRole.GetRole(tokenRole)
		if !ok {
			return respFail
		}

		requiredRoleLevel := u.GetRoleLevel(role.RoleId)

		if (tokenRoleLevel == u.Superuser && requiredRoleLevel == u.User) || (tokenRoleLevel == u.Admin && requiredRoleLevel > u.Admin && requiredRoleLevel <= u.Member) || (tokenRoleLevel == u.MerchantAdmin && requiredRoleLevel > u.MerchantAdmin && requiredRoleLevel <= u.Member) || (tokenRoleLevel == u.Root && requiredRoleLevel > u.Root && requiredRoleLevel <= u.Member) {

			return role.Delete(userId) //Update account
		} else {
			return u.Message(false, "0x11020:You do not have access authority")
		}
	} else {
		return u.Message(false, "0x11020:You do not have access authority")
	}
}
