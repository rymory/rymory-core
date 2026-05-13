// Copyright (c) 2017-2026 Onur Yaşar
// Licensed under AGPL v3 + Commercial Exception
// See LICENSE.txt

package main

import (
	"zombie"

	u "github.com/lemoras/goutils/api"
)

func Main(in zombie.Request) (*u.Response, error) {
	return zombie.Invoke(in)
}
