// Copyright (c) 2017-2026 Onur Yaşar
// Licensed under AGPL v3 + Commercial Exception
// See LICENSE.txt

package main

import (
	"member"

	u "github.com/lemoras/goutils/api"
)

func Main(in member.Request) (*u.Response, error) {
	return member.Invoke(in)
}
