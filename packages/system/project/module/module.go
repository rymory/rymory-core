package project

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	u "github.com/lemoras/goutils/api"
)

type Request struct {
	Http CustomHttp `json:"http"`

	AppId          int       `json:"appId"`
	MerchantId     uuid.UUID `json:"merchantId"`
	RequiredStatus int       `json:"requiredStatus"`
	ProjectName    string    `json:"projectName"`
	ManageMember   int       `json:"manageMember"`
	MemberStatus   int       `json:"memberStatus"`

	IsOwnerLemoras bool     `json:"isOwnerLemoras"`
	Domains        []string `json:"domains"`
}

type CustomHttp struct {
	CustomHeader CustomHeader `json:"headers"`
	Method       string       `json:"method"`
}

type CustomHeader struct {
	Authorization string `json:"authorization"`
	Referer       string `json:"referer"`
}

func Invoke(in Request) (*u.Response, error) {

	var resp map[string]interface{}

	if in.Http.CustomHeader.Authorization == "" {
		refererDomain := ""
		if strings.Contains(in.Http.CustomHeader.Referer, "http://") {
			refererDomain = strings.Split(in.Http.CustomHeader.Referer, "http://")[1]
		}
		if strings.Contains(in.Http.CustomHeader.Referer, "https://") {
			refererDomain = strings.Split(in.Http.CustomHeader.Referer, "https://")[1]
		}

		domain := strings.Split((refererDomain), "/")[0]
		return u.Respond(GetMerchantByDomain(domain))
	}

	context := &u.Context{}
	if isOk, res := u.JwtAuthentication(in.Http.CustomHeader.Authorization, context); !isOk {
		return &res, nil
	}

	switch in.Http.Method {
	case "GET":
		resp = GetProjects(*context)
		break
	case "POST":
		project := Project{}
		project.AppId = in.AppId
		project.MerchantId = in.MerchantId
		project.ProjectName = in.ProjectName
		if in.ManageMember == 1 {
			project.ManageMember = false
		}
		if in.ManageMember == 2 {
			project.ManageMember = true
		}
		domains := ""
		if len(in.Domains) > 0 {
			domains = strings.Join(in.Domains, ",")
		}
		project.Domains = domains
		project.IsOwnerLemoras = in.IsOwnerLemoras

		resp = CreateProject(project, *context)
		break
	case "PUT":
		requiredProject := RequiredProject{}
		requiredProject.AppId = in.AppId
		requiredProject.MerchantId = in.MerchantId
		requiredProject.RequiredStatus = in.RequiredStatus
		requiredProject.ProjectName = in.ProjectName
		requiredProject.ManageMember = in.ManageMember
		requiredProject.MemberStatus = in.MemberStatus

		resp = UpdateProject(requiredProject, *context)
		break
	case "DELETE":
		requiredProject := RequiredProject{}
		requiredProject.AppId = in.AppId
		requiredProject.MerchantId = in.MerchantId

		resp = DeleteProject(requiredProject, *context)
		break
	default:
		resp = u.Message(false, "0x11028:Invalid request")
		break
	}

	return u.Respond(resp)
}

var CreateProject = func(project Project, context u.Context) map[string]interface{} {

	if project.AppId == 0 {
		return u.Message(false, "0x11018:AppId is zero value")
	}

	if project.MerchantId == uuid.Nil {
		return u.Message(false, "0x11019:MerchantId is zero value")
	}

	if project.Domains == "" {
		return u.Message(false, "0x11037:Domains is empty")
	}

	if project.ProjectName == "" {
		return u.Message(false, "0x11038:Project name is empty")
	}

	roleId := context.RoleId
	tokenRoleLevel := u.GetRoleLevel(roleId)

	if tokenRoleLevel == u.Root { // || tokenRoleLevel == u.MerchantToken
		userId := context.UserId
		return project.Create(userId) //Create account
	}

	return u.Message(false, "0x11020:You do not have access authority")
}

var UpdateProject = func(requiredProject RequiredProject, context u.Context) map[string]interface{} {

	if requiredProject.AppId == 0 {
		return u.Message(false, "0x11018:AppId is zero value")
	}

	if requiredProject.MerchantId == uuid.Nil {
		return u.Message(false, "0x11019:MerchantId is zero value")
	}

	roleId := context.RoleId
	tokenRoleLevel := u.GetRoleLevel(roleId)

	merchantId := context.MerchantId
	appId := context.AppId

	if resp, ok := GetRole(context.UserId, context.AppId, context.MerchantId, context.RoleId); !ok {
		return resp
	}

	if (tokenRoleLevel == u.MerchantAdmin && merchantId == requiredProject.MerchantId && appId == requiredProject.AppId) || tokenRoleLevel == u.Root {

		project, ok := requiredProject.GetProject()

		if !ok {
			return u.Message(false, "0x11039:Project not found or already soft deleted")
		}

		userId := context.UserId

		if (project.Active && tokenRoleLevel == u.MerchantAdmin) || tokenRoleLevel == u.Root {
			if respD, result, okD := requiredProject.Update(userId); !okD {
				return respD
			} else {
				if requiredProject.ProjectName != "" {
					projectInfo := fmt.Sprint(result.ID) + ":" + result.ProjectName
					if respPIU, ok := BulkRoleUpdateProjectInfo(projectInfo, result.AppId, result.MerchantId); !ok {
						return respPIU
					}
				}
				if !result.Active && requiredProject.RequiredStatus != 0 {
					if _, okBD := BulkRoleDelete(project.AppId, project.MerchantId, true); !okBD {
						if resp, _, ok := requiredProject.Update(userId); !ok {
							return resp
						}
						return u.Message(false, "0x11040:Project passived but Users werent passived")
					}
				} else {
					if tokenRoleLevel == u.Root && requiredProject.RequiredStatus != 0 {
						if _, okBD := BulkRoleRevertDelete(project.AppId, project.MerchantId); !okBD {
							if resp, _, ok := requiredProject.Update(userId); !ok {
								return resp
							}
							return u.Message(false, "0x11041:Project activated but Users werent activated")
						}
					}
				}

				if result.Active && requiredProject.RequiredStatus != 1 {
					if requiredProject.ManageMember == 2 {
						requiredProject.MemberStatus = 2
					}
					BulkRoleMember(project.AppId, project.MerchantId, requiredProject.MemberStatus)
				}
			}
		} else {
			return u.Message(false, "0x11042:Project and users have passived. You can activate it through the portal from which you receive service")
		}

		return u.Message(true, "0x11043:Project and users updated by successful")
	}

	return u.Message(false, "0x11020:You do not have access authority")
}

var GetProjects = func(context u.Context) map[string]interface{} {

	roleId := context.RoleId
	tokenRoleLevel := u.GetRoleLevel(roleId)
	merchantId := context.MerchantId

	if tokenRoleLevel == u.Root || tokenRoleLevel == u.MerchantAdmin {
		resp, _ := GetProjectsByMerchantId(merchantId)
		return resp
	}

	return u.Message(false, "0x11020:You do not have access authority")
}

var GetMerchantByDomain = func(domain string) map[string]interface{} {

	resp, _ := GetProjectByDomain(domain)
	return resp
}

var DeleteProject = func(requiredProject RequiredProject, context u.Context) map[string]interface{} {

	if requiredProject.AppId == 0 {
		return u.Message(false, "0x11018:AppId is zero value")
	}

	if requiredProject.MerchantId == uuid.Nil {
		return u.Message(false, "0x11019:MerchantId is zero value")
	}

	roleId := context.RoleId
	tokenRoleLevel := u.GetRoleLevel(roleId)

	merchantId := context.MerchantId
	appId := context.AppId

	if resp, ok := GetRole(context.UserId, context.AppId, context.MerchantId, context.RoleId); !ok {
		return resp
	}

	if (tokenRoleLevel == u.MerchantAdmin && merchantId == requiredProject.MerchantId && appId == requiredProject.AppId) || tokenRoleLevel == u.Root {

		project, ok := requiredProject.GetProject()

		if !ok {
			return u.Message(false, "0x11039:Project not found or already soft deleted")
		}

		userId := context.UserId

		if tokenRoleLevel == u.Root {
			appId = requiredProject.AppId
		}

		if respD, okD := project.Delete(userId); !okD {
			return respD
		}

		if _, okBD := BulkRoleDelete(appId, merchantId, false); !okBD {
			if resp, ok := project.RevertDelete(); !ok {
				return resp
			}
			return u.Message(false, "0x11044:Project deleted but Users werent deleted")
		}

		return u.Message(true, "0x11045:Project deleted and Users deleted by successful")
	}

	return u.Message(false, "0x11020:You do not have access authority")
}
