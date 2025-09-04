package authenticate

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	u "github.com/lemoras/goutils/api"
)

type CookieFeature struct {
	Secure   bool
	SameSite SameSite
	Domain   string
}

type SameSite int

const (
	SameSiteLax SameSite = iota
	SameSiteStrict
)

var SetJWTAutCookie = func(httpToken string, requestOrigin string, cookie CookieFeature) map[string]string {

	var headers map[string]string
	if httpToken != "" {
		token := httpToken

		// Cookie stringini hazırla
		exp := time.Now().Add(20 * time.Minute).UTC().Format(time.RFC1123)

		var cookieStr string

		cookieStr = fmt.Sprintf(
			"authToken=Bearer %s; Expires=%s; Path=/; HttpOnly;",
			token, exp,
		)

		if cookie.Secure {
			cookieStr = cookieStr + " Secure;"
		}

		if cookie.Domain != "" {
			cookieStr = cookieStr + fmt.Sprintf(" Domain=%s;", cookie.Domain)
		}

		switch cookie.SameSite {
		case SameSiteLax:
			cookieStr = cookieStr + " SameSite=Lax;"
		case SameSiteStrict:
			cookieStr = cookieStr + " SameSite=Strict;"
		default:
			cookieStr = cookieStr + " SameSite=None;"
		}

		if !isOriginAllowed(requestOrigin) {
			return headers
		}

		headers = map[string]string{
			"Set-Cookie":                       cookieStr,
			"Access-Control-Allow-Origin":      requestOrigin,
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Headers":     "Content-Type, Authorization, Cookie",
			"Access-Control-Allow-Methods":     "GET, POST, OPTIONS",
			"Content-Type":                     "application/json",
		}
	}

	return headers
}

var CheckJWTAutCookie = func(requestToken string, context *u.Context, headers u.CustomHeader) (bool, u.Response) {

	if headers.XAPIKey == os.Getenv("x_api_key") {
		return u.JwtAuthentication(requestToken, context)
	}

	if headers.Cookie == "" {
		return u.ResMessage(false, "Missing Cookie")
	}

	cookies := strings.Split(headers.Cookie, "; ")

	var authTokenValue string

	for _, cookie := range cookies {
		parts := strings.SplitN(cookie, "=", 2)
		if len(parts) == 2 && parts[0] == "authToken" {
			authTokenValue = parts[1]
			break
		}
	}

	if authTokenValue == "" {
		return u.ResMessage(false, "0x11130:Missing auth token")
	}

	return u.JwtAuthentication(authTokenValue, context)
}

var CheckAuthEmpty = func(headers u.CustomHeader) bool {

	if headers.XAPIKey == os.Getenv("x_api_key") {
		return headers.Authorization == ""
	}

	if headers.Cookie == "" {
		return true
	}

	tokenValue := ""
	cookies := strings.Split(headers.Cookie, "; ")
	for _, cookie := range cookies {
		parts := strings.Split(cookie, "=")
		if len(parts) == 2 && parts[0] == "authToken" {
			tokenValue = parts[1]
			break
		}
	}

	return tokenValue == ""
}

func isOriginAllowed(origin string) bool {

	allowedDomainsStr := os.Getenv("cookie_allowed_domains")

	allowedDomains := strings.Split(allowedDomainsStr, ",")

	parsedOrigin, err := url.Parse(origin)
	if err != nil {
		return false
	}

	originHost := parsedOrigin.Host

	if strings.HasPrefix(originHost, "dev.local:") {
		originHost = "dev.local"
	}

	for _, domain := range allowedDomains {
		if originHost == domain {
			return true
		}

		if strings.HasSuffix(originHost, "."+domain) {
			return true
		}
	}

	return false
}
