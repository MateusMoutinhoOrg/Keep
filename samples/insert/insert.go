package main

import (
	keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/native"
	keep_lib "github.com/MateusMoutinhoOrg/Keep/pkg/keep"
)

func main() {
	deps := keep_deps.New()
	keep := keep_lib.New(deps)

	mateus := keep.Insert("user")

	mateus.SetPrimaryKey("email", "mateusmoutinho01@gmail.com")
	mateus.SetPrimaryKey("username", "mateus")
	mateus.AddString("password", "12345678")

}
