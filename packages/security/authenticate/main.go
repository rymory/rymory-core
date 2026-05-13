// Copyright (c) 2017-2026 Onur Yaşar
// Licensed under AGPL v3 + Commercial Exception
// See LICENSE.txt

// https://github.com/rymory/rymory-core
// rymory.org 
// onuryasar.org
// onxorg@proton.me 

package main

import (
	"authenticate"

	u "github.com/rymory/goutils/api"
)

func Main(in authenticate.Request) (*u.Response, error) {
	return authenticate.Invoke(in)
}
