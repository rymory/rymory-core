package main

import (
	"authenticate"

	u "gitlab.com/onxorg/goutils/api"
)

func Main(in authenticate.Request) (*u.Response, error) {
	return authenticate.Invoke(in)
}
