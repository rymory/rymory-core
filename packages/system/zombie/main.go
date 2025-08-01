package main

import (
	"zombie"

	u "github.com/lemoras/goutils/api"
)

func Main(in zombie.Request) (*u.Response, error) {
	return zombie.Invoke(in)
}
