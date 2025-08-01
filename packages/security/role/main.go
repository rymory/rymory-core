package main

import (
	"role"

	u "github.com/lemoras/goutils/api"
)

func Main(in role.Request) (*u.Response, error) {
	return role.Invoke(in)
}
