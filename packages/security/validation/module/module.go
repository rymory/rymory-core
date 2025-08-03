package validation

import (
	"encoding/json"
	"fmt"
	"net/http"

	u "github.com/lemoras/goutils/api"
)

type Request struct {
	Http CustomHttp `json:"http"`
}

type CustomHttp struct {
	CustomHeader CustomHeader `json:"headers"`
}

type CustomHeader struct {
	Authorization string `json:"authorization"`
}

func Invoke(in Request) (*u.Response, error) {

	context := &u.Context{}
	if isOk, res := u.JwtAuthentication(in.Http.CustomHeader.Authorization, context); !isOk {
		return &res, nil
	}

	headData := make(map[string]string)
	headData["UserId"] = fmt.Sprint(context.UserId)
	headData["RoleId"] = fmt.Sprint(context.RoleId)
	headData["AppId"] = fmt.Sprint(context.AppId)
	headData["MerchantId"] = fmt.Sprint(context.MerchantId)
	headData["HasId"] = fmt.Sprint(context.HasId)
	headData["ProjectId"] = fmt.Sprint(context.ProjectId)
	headData["CustomData"] = fmt.Sprint(context.CustomData)
	headData["InitCompleted"] = fmt.Sprint(context.InitCompleted)

	headData["Access-Control-Expose-Headers"] = "userid"
	headData["Access-Control-Expose-Headers"] = "UserId"

	// headData["Access-Control-Expose-Headers"] = "roleId"
	// headData["Access-Control-Expose-Headers"] = "RoleId"

	// headData["Access-Control-Expose-Headers"] = "appId"
	// headData["Access-Control-Expose-Headers"] = "AppId"

	// headData["Access-Control-Expose-Headers"] = "merchantId"
	// headData["Access-Control-Expose-Headers"] = "MerchantId"

	// headData["Access-Control-Expose-Headers"] = "hasId"
	// headData["Access-Control-Expose-Headers"] = "HasId"

	// headData["Access-Control-Expose-Headers"] = "projectId"
	// headData["Access-Control-Expose-Headers"] = "ProjectId"

	// headData["Access-Control-Expose-Headers"] = "customData"
	// headData["Access-Control-Expose-Headers"] = "CustomData"

	jsonData, _ := json.Marshal(u.Message(true, "0x11031:Success"))

	return &u.Response{
		StatusCode: http.StatusOK,
		Headers:    headData,
		Body:       string(jsonData),
	}, nil
}
