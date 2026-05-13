// Copyright (c) 2017-2026 Onur Yaşar
// Licensed under AGPL v3 + Commercial Exception
// See LICENSE.txt

package main

import (
	"role"

	u "github.com/lemoras/goutils/api"
)

func Main(in role.Request) (*u.Response, error) {
	return role.Invoke(in)
}
