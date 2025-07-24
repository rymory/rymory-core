package main

import (
	"project"

	u "gitlab.com/onxorg/goutils/api"
)

func Main(in project.Request) (*u.Response, error) {
	return project.Invoke(in)
}
