package main

import (
	"initialize"

	u "gitlab.com/onxorg/goutils/api"
)

func Main(in initialize.Request) (*u.Response, error) {
	return initialize.Invoke(in)
}
