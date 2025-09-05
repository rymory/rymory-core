package authenticate

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	u "github.com/lemoras/goutils/api"
)

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`

	RoleId     int       `json:"roleId"`
	AppId      int       `json:"appId"`
	MerchantId uuid.UUID `json:"merchantId"`
	HasId      bool      `json:"hasId"`
	CustomData string    `json:"customData"`

	UserId uuid.UUID `json:"userId"`

	Token string `json:"token"`

	Http u.CustomHttp `json:"http"`
}

func Invoke(in Request) (*u.Response, error) {

	var resp map[string]interface{}
	var httpToken string
	var headers map[string]string
	var cookie CookieFeature

	isHttpOnlyAuthCookieStr := os.Getenv("isHttpOnlyAuthCookie")

	isHttpOnlyAuthCookie, err := strconv.ParseBool(isHttpOnlyAuthCookieStr)
	if err != nil {
		// Geçersiz değer olursa default false
		log.Printf("invalid value for isHttpOnlyAuthCookie: %s (defaulting to false)", isHttpOnlyAuthCookieStr)
		isHttpOnlyAuthCookie = false
	}

	if CheckAuthEmpty(in.Http.CustomHeader) {
		account := &Account{}
		account.Email = strings.ToLower(in.Email)
		account.Password = in.Password
		resp, httpToken = Login(account.Email, account.Password, uuid.Nil)
		if isHttpOnlyAuthCookie {
			isHttpSSLStr := os.Getenv("isHttpSSL")

			isHttpSSL, err := strconv.ParseBool(isHttpSSLStr)
			if err != nil {
				// Geçersiz değer olursa default false
				log.Printf("invalid value for isHttpOnlyAuthCookie: %s (defaulting to false)", isHttpOnlyAuthCookieStr)
				isHttpSSL = false
			}
			cookie = CookieFeature{
				Secure:     isHttpSSL,
				Domain:     ".dev.local",
				SameSite:   SameSiteLax,
				CookieName: "low_authToken",
			}
		}
	} else {
		context := &u.Context{}
		if isOk, res := CheckJWTAutCookie(in.Http.CustomHeader.Authorization, context, in.Http.CustomHeader); !isOk {
			return &res, nil
		}
		fmt.Println("Burda")
		if isHttpOnlyAuthCookie {
			if in.Http.Method == "GET" {
				path := strings.Replace(in.Http.Path, "/", "", -1)
				fmt.Println("path")
				if path != "" {
					err := json.Unmarshal([]byte(path), &in)
					if err != nil {
						fmt.Println("json parsing error %s", err.Error())
					}
					in.Http.Path = ""
				}
			}
			fmt.Println(json.Marshal(&in))
			fmt.Println("buralari duzgun gecti")
			cookieDomain := strings.TrimPrefix(in.Http.CustomHeader.Origin, "http://")
			cookieDomain = strings.TrimPrefix(cookieDomain, "https://")
			fmt.Println(cookieDomain)
			isHttpSSLStr := os.Getenv("isHttpSSL")

			isHttpSSL, err := strconv.ParseBool(isHttpSSLStr)
			if err != nil {
				// Geçersiz değer olursa default false
				log.Printf("invalid value for isHttpOnlyAuthCookie: %s (defaulting to false)", isHttpOnlyAuthCookieStr)
				isHttpSSL = false
			}
			cookie = CookieFeature{
				Secure:     isHttpSSL,
				Domain:     "",
				SameSite:   SameSiteLax,
				CookieName: "strong_authToken",
			}
		}

		if false {
			resp, httpToken = LoginByToken(in.Email, in.Token, *context)
		} else {

			member := Membership{}

			member.AppId = in.AppId
			member.RoleId = in.RoleId
			member.MerchantId = in.MerchantId
			member.HasId = in.HasId
			member.CustomData = in.CustomData

			path := strings.Replace(in.Http.Path, "/", "", -1)
			if path == "securityauthenticate" {
				path = ""
			}

			switch path {
			case "":
				resp, httpToken = RenewToken(member, *context)
				fmt.Println(httpToken)
			case "change":
				member.UserId = in.UserId
				resp, httpToken = ChangeSession(member, *context)
			}
		}
	}

	if isHttpOnlyAuthCookie {
		headers = SetJWTAutCookie(httpToken, in.Http.CustomHeader.Origin, cookie)
		return u.RespondWithHeaders(resp, headers)
	}

	return u.Respond(resp)
}

var LoginByToken = func(email string, token string, context u.Context) (map[string]interface{}, string) {

	appId := context.AppId
	merchantId := context.MerchantId
	roleId := context.RoleId
	tokenRoleLevel := u.GetRoleLevel(roleId)

	if appId == 0 && merchantId == uuid.Nil && tokenRoleLevel == u.None {
		id := context.UserId
		resp, httpToken := Login(email, "", id)
		return resp, httpToken
	}

	return u.Message(false, "0x11017:It doesnt build a new token by strong token"), ""
}

var RenewToken = func(member Membership, context u.Context) (map[string]interface{}, string) {

	appId := context.AppId
	merchantId := context.MerchantId
	roleId := context.RoleId
	tokenRoleLevel := u.GetRoleLevel(roleId)

	if appId == 0 && merchantId == uuid.Nil && tokenRoleLevel == u.None {
		id := context.UserId
		resp, httpToken := BuildToken(id, member.RoleId, member.AppId, member.MerchantId, false, member.CustomData)
		return resp, httpToken
	}

	return u.Message(false, "0x11017:It doesnt build a new token by strong token"), ""
}

var ChangeSession = func(member Membership, context u.Context) (map[string]interface{}, string) {

	tokenRoleLevel := u.GetRoleLevel(context.RoleId)

	if tokenRoleLevel == u.Root {

		if member.AppId == 0 {
			return u.Message(false, "0x11018:AppId is zero value"), ""
		}

		if member.MerchantId == uuid.Nil {
			return u.Message(false, "0x11019:MerchantId is zero value"), ""
		}
	} else if tokenRoleLevel == u.MerchantAdmin {
		if member.AppId == 0 {
			member.AppId = context.AppId
		}
		member.MerchantId = context.MerchantId
	} else if tokenRoleLevel > u.MerchantAdmin && tokenRoleLevel < u.User {
		member.AppId = context.AppId
		member.MerchantId = context.MerchantId
	} else {
		return u.Message(false, "0x11020:You do not have access authority"), ""
	}

	if (context.AppId == member.AppId && context.MerchantId == member.MerchantId) || (tokenRoleLevel == u.MerchantAdmin && context.MerchantId == member.MerchantId) || (tokenRoleLevel == u.Root) {

		requiredRoleLevel := u.GetRoleLevel(member.RoleId)

		if (tokenRoleLevel == u.Superuser && requiredRoleLevel == u.User) || (tokenRoleLevel == u.Admin && requiredRoleLevel > u.Admin && requiredRoleLevel < u.Member) || (tokenRoleLevel == u.MerchantAdmin && requiredRoleLevel > u.MerchantAdmin && requiredRoleLevel < u.Member) || (tokenRoleLevel == u.Root && requiredRoleLevel > u.Root && requiredRoleLevel < u.Member) {
			if u.CheckOk(CheckUser(context.UserId, context.RoleId, context.AppId, context.MerchantId)) {
				return BuildToken(member.UserId, member.RoleId, member.AppId, member.MerchantId, member.UserId != context.UserId, member.CustomData)
			}
			return u.Message(false, "0x11020:You do not have access authority"), ""
		} else {
			return u.Message(false, "0x11020:You do not have access authority"), ""
		}
	} else {
		return u.Message(false, "0x11020:You do not have access authority"), ""
	}
}
