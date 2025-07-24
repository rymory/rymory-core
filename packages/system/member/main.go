package main

import (
	"member"

	u "gitlab.com/onxorg/goutils/api"
)

func Main(in member.Request) (*u.Response, error) {
	return member.Invoke(in)
}
