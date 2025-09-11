package account

import (
	"strings"

	u "github.com/lemoras/goutils/api"
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

	var resp map[string]interface{}

	switch in.Http.Method {
	case "GET":
		context := &u.Context{}
		if isOk, res := u.JwtAuthentication(in.Http.CustomHeader.Authorization, context); !isOk {
			return &res, nil
		}
		resp = GetAccount(*context)
	case "POST":
		account := Account{}

		account.Password = in.Password
		account.AboutMe = in.AboutMe
		account.PhotoUrl = in.PhotoUrl
		account.Nickname = in.Nickname
		account.Email = strings.ToLower(in.Email)

		resp = CreateAccount(account)

	case "PUT":

		context := &u.Context{}
		if isOk, res := u.JwtAuthentication(in.Http.CustomHeader.Authorization, context); !isOk {
			return &res, nil
		}
		account := Account{}

		account.Password = in.Password
		account.AboutMe = in.AboutMe
		account.PhotoUrl = in.PhotoUrl
		account.Nickname = in.Nickname

		resp = UpdateAccount(account, *context)
	}
	return u.Respond(resp)
}

var CreateAccount = func(account Account) map[string]interface{} {
	return account.Create() //Create account
}

var UpdateAccount = func(account Account, context u.Context) map[string]interface{} {

	return account.Update(context.UserId, context.HasId) //Update account
}

var GetAccount = func(context u.Context) map[string]interface{} {

	return Get(context.UserId) //Get account
}
