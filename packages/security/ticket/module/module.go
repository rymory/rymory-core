package ticket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	u "github.com/lemoras/goutils/api"
)

type Request struct {
	CustomData string `json:"customData"`

	Http u.CustomHttp `json:"http"`
}

func Invoke(in Request) (*u.Response, error) {

	if in.Http.CustomHeader.Authorization != "" {

		context := &u.Context{}
		if isOk, res := u.JwtAuthentication(in.Http.CustomHeader.Authorization, context); !isOk {
			return &res, nil
		}

		if in.CustomData != "" {

			if in.Http.Method == "POST" {
				return u.Respond(GenerateTicket(in.CustomData, *context))
			} else if in.Http.Method == "GET" {

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

				resp := ValidTicket(in.CustomData, *context)

				jsonData, _ := json.Marshal(resp)

				return &u.Response{
					StatusCode: http.StatusOK,
					Headers:    headData,
					Body:       string(jsonData),
				}, nil
			}
		}

	}

	return u.Respond(u.Message(false, "Invlaid"))
}

var GenerateTicket = func(customData string, context u.Context) map[string]interface{} {

	appId := context.AppId
	merchantId := context.MerchantId
	roleId := context.RoleId
	tokenRoleLevel := u.GetRoleLevel(roleId)

	if appId > 0 && merchantId != uuid.Nil && tokenRoleLevel != u.None {
		id := context.UserId
		resp := BuildToken(id, roleId, appId, merchantId, customData)
		return resp
	}

	return u.Message(false, "0x11017:It doesnt build a new token by strong token")
}

var ValidTicket = func(customData string, context u.Context) map[string]interface{} {

	ticketContext := &u.Context{}

	if isOk, res := JwtTicket(customData, ticketContext); !isOk {
		return res
	}

	if ticketContext.UserId == context.UserId &&
		ticketContext.MerchantId == context.MerchantId &&
		ticketContext.AppId == context.AppId &&
		ticketContext.RoleId == context.RoleId {

		resp := u.Message(true, "Ticket is valid.")
		resp["ticket"] = ticketContext.CustomData

		return resp
	}

	return u.Message(false, "Ticket is invalid.")
}

var JwtTicket = func(requestToken string, ticketContext *u.Context) (bool, map[string]interface{}) {

	if requestToken == "" { //Token is missing, returns with error code 403 Unauthorized
		return false, u.Message(false, "0x11130:Missing auth ticket")
	}

	splitted := strings.Split(requestToken, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	if len(splitted) != 2 {
		return false, u.Message(false, "0x11143:Invalid/Malformed auth ticket")
	}

	tokenPart := splitted[1] //Grab the token part, what we are truly interested in
	tk := &u.Token{}

	token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ticket_secret_key")), nil
	})

	if err != nil {
		return false, u.Message(false, "0x11144:Malformed authentication token")
	}

	if !token.Valid { //Token is invalid, maybe not signed on this server
		return false, u.Message(false, "Ticket is not valid.")
	}

	ticketContext.UserId = tk.UserId
	ticketContext.RoleId = tk.RoleId
	ticketContext.AppId = tk.AppId
	ticketContext.MerchantId = tk.MerchantId
	ticketContext.HasId = tk.HasId
	ticketContext.ProjectId = tk.ProjectId
	ticketContext.CustomData = tk.CustomData
	ticketContext.InitCompleted = tk.InitCompleted

	return true, u.Message(true, "Ticket is valid.")
}
