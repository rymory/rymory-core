package main

import (
	"member"

	u "github.com/lemoras/goutils/api"
)

func Main(in member.Request) (*u.Response, error) {
	return member.Invoke(in)
}
