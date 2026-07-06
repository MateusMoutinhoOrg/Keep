package main

import (
	keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/native"
	keep_lib "github.com/MateusMoutinhoOrg/Keep/pkg/keep"
)

func main() {
	deps := keep_deps.New()
	keep := keep_lib.New(deps)

}
