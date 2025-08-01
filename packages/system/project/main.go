package main

import (
	"project"

	u "github.com/lemoras/goutils/api"
)

func Main(in project.Request) (*u.Response, error) {
	return project.Invoke(in)
}
