package ticket

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	u "github.com/lemoras/goutils/api"
)

func BuildToken(userId uuid.UUID, roleId int, appId int, merchantId uuid.UUID, customData string) map[string]interface{} {

	atClaims := jwt.StandardClaims{}
	atClaims.Issuer = os.Getenv("JWT_ISSUER")
	atClaims.ExpiresAt = time.Now().Add(time.Minute * 1).Unix()

	//Create JWT token
	tk := &u.Token{UserId: userId, RoleId: roleId, AppId: appId, MerchantId: merchantId, CustomData: customData, StandardClaims: atClaims}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS512"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("TICKET_SECRET_KEY")))

	resp := u.Message(true, "0x11023:Ticket done success")
	resp["ticket"] = tokenString
	return resp
}
