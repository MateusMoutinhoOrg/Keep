package main

import (
	"fmt"

	keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/native"
	"github.com/MateusMoutinhoOrg/Keep/pkg/database"
	keep_lib "github.com/MateusMoutinhoOrg/Keep/pkg/keep"
)

func main() {
	deps := keep_deps.New()
	keep := keep_lib.New(deps)

	database := keep.NewDatabase(database.Schema{})
	fmt.Println(database)
}
