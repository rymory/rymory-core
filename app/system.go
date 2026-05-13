// Copyright (c) 2017-2026 Onur Yaşar
// Licensed under AGPL v3 + Commercial Exception
// See LICENSE.txt

// https://github.com/rymory/rymory-core
// rymory.org 
// onuryasar.org
// onxorg@proton.me 

package main

import (
	"encoding/json"
	"fmt"
	"initialize"
	"io/ioutil"
	"member"
	"net/http"
	"project"
	"strings"
	"zombie"
)

func Initialize(w http.ResponseWriter, r *http.Request) {

	var in initialize.Request

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &in)

	fmt.Println("Endpoint hit: Initialize")

	in.Http.CustomHeader.Authorization = r.Header.Get("authorization")
	in.Http.Path = strings.Replace(r.URL.Path, "system/init/", "", -1)
	in.Http.Method = r.Method

	fmt.Println(in.Http.Path)

	resp, _ := initialize.Invoke(in)
	w.Write([]byte(resp.Body))
}

func Member(w http.ResponseWriter, r *http.Request) {

	var in member.Request

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &in)

	fmt.Println("Endpoint hit: Member")

	in.Http.CustomHeader.Authorization = r.Header.Get("authorization")
	in.Http.Method = r.Method

	resp, _ := member.Invoke(in)
	w.Write([]byte(resp.Body))
}

func Project(w http.ResponseWriter, r *http.Request) {

	var in project.Request

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &in)

	fmt.Println("Endpoint hit: Project")

	in.Http.CustomHeader.Authorization = r.Header.Get("authorization")
	in.Http.Method = r.Method
	in.Http.CustomHeader.Referer = r.Header.Get("referer")

	resp, _ := project.Invoke(in)
	w.Write([]byte(resp.Body))
}

func Zombie(w http.ResponseWriter, r *http.Request) {

	var in zombie.Request

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &in)

	fmt.Println("Endpoint hit: Zombie")

	in.Http.CustomHeader.Authorization = r.Header.Get("authorization")

	resp, _ := zombie.Invoke(in)
	w.Write([]byte(resp.Body))
}
