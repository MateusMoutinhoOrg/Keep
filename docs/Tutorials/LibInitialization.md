# Library Initialization

## Description
Covers installing the library and initializing it with the standard (filesystem) adapter in a new program. To pick a different backend or write your own, see [DepsMechanic.md](/docs/Explanation/DepsMechanic.md).

### Rules
- Requires Go 1.22 or newer.

---

## Workflow
1. Install the lib:
   ```bash
   go get github.com/MateusMoutinhoOrg/Keep@v0.0.1
   ```
2. Create a file called `main.go` with the following code:
   ```go
   package main

   // 1. Import an adapter, the schema types, and the lib
   import (
       "fmt"

       keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
       "github.com/MateusMoutinhoOrg/Keep/pkg/database"
       keep_lib "github.com/MateusMoutinhoOrg/Keep/pkg/keep"
   )

   // 2. Describe your data: one "user" collection with three fields
   var Props = database.Props{
       Path: "myDatabase/",
       Schemas: []database.Schema{
           {
               Name: "user",
               Itens: []database.Item{
                   {Name: "email", Type: database.Key, Required: true},
                   {Name: "username", Type: database.Key, Required: true},
                   {Name: "age", Type: database.Int, Required: true},
               },
           },
       },
   }

   func main() {
       // 3. Create deps via an adapter (the "opinionated" part)
       deps := keep_deps.New()

       // 4. Inject deps into the pure library
       keep := keep_lib.New(deps)

       // 5. Use the library — it never knows which adapter is behind the scenes
       db := keep.NewDatabase(Props)
       users := db.GetSchema("user")

       created, err := users.NewItem(map[string]any{
           "email":    "mateus@gmail.com",
           "username": "mateus",
           "age":      27,
       })
       if err != nil {
           fmt.Println("error creating user:", err)
           return
       }
       fmt.Println("created:", created)
   }
   ```
3. Run the code:
   ```bash
   go run main.go
   ```
4. Describe the rest of your data with [Schemas.md](/docs/Explanation/Schemas.md), and operate on it with [Records.md](/docs/Explanation/Records.md).
