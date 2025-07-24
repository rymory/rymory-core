package main

import (
	"zombie"

	u "gitlab.com/onxorg/goutils/api"
)

func Main(in zombie.Request) (*u.Response, error) {
	return zombie.Invoke(in)
}
