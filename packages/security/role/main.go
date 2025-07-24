package main

import (
	"role"

	u "gitlab.com/onxorg/goutils/api"
)

func Main(in role.Request) (*u.Response, error) {
	return role.Invoke(in)
}
