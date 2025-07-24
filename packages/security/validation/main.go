package main

import (
	"validation"

	u "gitlab.com/onxorg/goutils/api"
)

func Main(in validation.Request) (*u.Response, error) {
	return validation.Invoke(in)
}
