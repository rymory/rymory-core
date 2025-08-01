package main

import (
	"authenticate"

	u "github.com/lemoras/goutils/api"
)

func Main(in authenticate.Request) (*u.Response, error) {
	return authenticate.Invoke(in)
}
