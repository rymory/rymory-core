package main

import (
	"validation"

	u "github.com/lemoras/goutils/api"
)

func Main(in validation.Request) (*u.Response, error) {
	return validation.Invoke(in)
}
