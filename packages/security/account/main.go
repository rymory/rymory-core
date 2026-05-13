// Copyright (c) 2017-2026 Onur Yaşar
// Licensed under AGPL v3 + Commercial Exception
// See LICENSE.txt

package main

import (
	"account"

	u "github.com/lemoras/goutils/api"
)

// Main forwarding to Hello
// func Main(args map[string]interface{}) map[string]interface{} {
// 	fmt.Println("Main")
// 	return account.Main(args)
// }

func Main(in account.Request) (*u.Response, error) {
	return account.Invoke(in)
}
