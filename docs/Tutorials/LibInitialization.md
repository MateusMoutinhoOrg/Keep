# Library Initialization

## Description
Covers installing the library and initializing it with the standard (filesystem) adapter in a new program. To pick a different backend or write your own, see [DepsMechanic.md](/docs/Explanations/DepsMechanic.md); the shipped ones are listed in [Adapters.md](/docs/References/Adapters.md).

### Rules
- Requires Go 1.22 or newer.
- The `sandbox` package is named `lib`, so import it under that alias: `lib "github.com/MateusMoutinhoOrg/Keep/sandbox"`.
- The schema description and the typed error come from `sandbox/contracts/api`.

---

## Workflow
1. Install the lib:
   ```bash
   go get github.com/MateusMoutinhoOrg/Keep@v0.0.2
   ```
2. Create a file called `main.go` with the following code:
   ```go
   package main

   // 1. Import an adapter, the lib, and the api contracts
   import (
       "fmt"

       "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
       lib "github.com/MateusMoutinhoOrg/Keep/sandbox"
       "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
   )

   // 2. Describe your data: one "user" collection with three fields
   func createProps() api.Props {

       //========================User==========================
       email := lib.NewKeyItem("email", true)
       username := lib.NewKeyItem("username", true)
       age := lib.NewIntItem("age", true)
       user := lib.NewSchema("user", email, username, age)

       //========================Props==========================
       return lib.NewProps("myDatabase/", user)
   }

   func main() {
       // 3. Create deps via an adapter (the "opinionated" part)
       deps := standard.New()

       // 4. Inject deps into the closed sandbox
       keep := lib.New(deps)

       // 5. Use the library — it never knows which adapter is behind the scenes
       props := createProps()
       db := keep.NewDatabase(props)
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
4. Describe the rest of your data with [Schemas.md](/docs/Explanations/Schemas.md), and operate on it with [Records.md](/docs/Explanations/Records.md).
