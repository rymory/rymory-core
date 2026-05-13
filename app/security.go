// Copyright (c) 2017-2026 Onur Yaşar
// Licensed under AGPL v3 + Commercial Exception
// See LICENSE.txt

package main

import (
	"account"
	"authenticate"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"role"
	"ticket"
	"validation"
)

func Authenticate(w http.ResponseWriter, r *http.Request) {

	var in authenticate.Request

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &in)

	fmt.Println("Endpoint hit: Authenticate")

	in.Http.CustomHeader.Authorization = r.Header.Get("Authorization")
	in.Http.Method = r.Method
	fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX   ")
	fmt.Println("XXXXXXXXXXXXXXXXXXXX:   " + r.Header.Get("Authorization"))

	fmt.Println("XXXXXXXXXXXXXXXXXXXX:   " + r.Header.Get("authorization"))

	resp, _ := authenticate.Invoke(in)
	w.Write([]byte(resp.Body))

	//json.NewEncoder(w).Encode(resp)
}

func Account(w http.ResponseWriter, r *http.Request) {

	var in account.Request

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &in)

	fmt.Println("Endpoint hit: Account")

	in.Http.CustomHeader.Authorization = r.Header.Get("Authorization")
	in.Http.Method = r.Method

	fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX   ")
	fmt.Println("XXXXXXXXXXXXXXXXXXXX:   " + r.Header.Get("Authorization"))

	fmt.Println("XXXXXXXXXXXXXXXXXXXX:   " + r.Header.Get("authorization"))
	resp, _ := account.Invoke(in)
	w.Write([]byte(resp.Body))
}

func Role(w http.ResponseWriter, r *http.Request) {

	var in role.Request

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &in)

	fmt.Println("Endpoint hit: Role")

	in.Http.CustomHeader.Authorization = r.Header.Get("authorization")
	in.Http.Method = r.Method

	resp, _ := role.Invoke(in)
	w.Write([]byte(resp.Body))
}

func Validation(w http.ResponseWriter, r *http.Request) {

	var in validation.Request

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &in)

	fmt.Println("Endpoint hit: Validation")

	in.Http.CustomHeader.Authorization = r.Header.Get("authorization")

	resp, _ := validation.Invoke(in)

	w.Header().Set("userId", resp.Headers["UserId"])
	w.Header().Set("roleId", resp.Headers["RoleId"])
	w.Header().Set("appId", resp.Headers["AppId"])
	w.Header().Set("merchantId", resp.Headers["MerchantId"])
	w.Header().Set("hasId", resp.Headers["HasId"])
	w.Header().Set("projectId", resp.Headers["ProjectId"])
	w.Header().Set("customData", resp.Headers["CustomData"])
	w.Header().Set("initCompleted", resp.Headers["InitCompleted"])

	// w.Header().Set("Access-Control-Expose-Headers", "userid")
	// w.Header().Set("Access-Control-Expose-Headers", "UserId")

	// w.ExposedHeaders([]string{"Access-Control-Expose-Headers", "userId", "UserId"})

	w.Write([]byte(resp.Body))
}

func Ticket(w http.ResponseWriter, r *http.Request) {

	var in ticket.Request

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &in)

	fmt.Println("Endpoint hit: Ticket")

	in.Http.CustomHeader.Authorization = r.Header.Get("authorization")
	in.Http.Method = r.Method

	resp, _ := ticket.Invoke(in)

	w.Header().Set("userId", resp.Headers["UserId"])
	w.Header().Set("roleId", resp.Headers["RoleId"])
	w.Header().Set("appId", resp.Headers["AppId"])
	w.Header().Set("merchantId", resp.Headers["MerchantId"])
	w.Header().Set("hasId", resp.Headers["HasId"])
	w.Header().Set("projectId", resp.Headers["ProjectId"])
	w.Header().Set("customData", resp.Headers["CustomData"])
	w.Header().Set("initCompleted", resp.Headers["InitCompleted"])

	w.Write([]byte(resp.Body))
}
