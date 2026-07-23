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

   // 1. Import an adapter and the lib
   import (
       "fmt"

       "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
       "github.com/MateusMoutinhoOrg/Keep/pkg/lib"
   )

   // 2. Describe your data: one "user" collection with three fields
   var Props = lib.Props{
       Path: "myDatabase/",
       Schemas: []lib.Schema{
           {
               Name: "user",
               Itens: []lib.Item{
                   {Name: "email", Type: lib.Key, Required: true},
                   {Name: "username", Type: lib.Key, Required: true},
                   {Name: "age", Type: lib.Int, Required: true},
               },
           },
       },
   }

   func main() {
       // 3. Create deps via an adapter (the "opinionated" part)
       deps := standard.New()

       // 4. Inject deps into the pure library
       keep := lib.New(deps)

       // 5. Use the library — it never knows which adapter is behind the scenes
       db := lib.NewDatabase(Props)
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
