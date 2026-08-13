# Struct Contracts

## Description
Explains why every contract in this project — `deps.Deps` and everything in `sandbox/contracts/api` — is a **struct of function fields** instead of an interface, what that buys, and what it costs.

---

## The Shape

A contract is a struct whose fields are functions. The library declares the shape; whoever fills it decides the behavior.

```go
// sandbox/contracts/deps/deps.go — what the library needs
type Deps struct {
	Write func(key string, value []byte) error
	Read  func(key string) ([]byte, error)
	// … one field per remaining requirement
}

// sandbox/contracts/api/api.go — what the library hands back
type SchemaInstance struct {
	Deps      deps.Deps
	Items     []Item
	Prefix    string
	NewItem   func(fields map[string]any) (SchemaItem, *Error)
	FindByKey func(key string, keyValue any) (SchemaItem, bool)
}
```

Callers use both exactly as they would an interface — `users.NewItem(fields)` reads the same whether `NewItem` is a method or a field holding a function. The difference only shows up at the wiring points.

An object that carries behavior — `Lib`, `KeepDatabase`, `SchemaInstance`, `SchemaItem` — leads with its own `Deps` field. That is what removes the need for a second, internal mirror type: the struct handed to the caller is the same struct the library's own code reads its dependencies from. `Item`, `Schema`, `Props`, and `Error` carry no behavior at all, so they have no `Deps` field either — they are plain data, built directly with a composite literal.

---

## Factories Fill the Fields

There is no internal type and no method set. `sandbox/lib/` holds **factories**: functions that take a pointer to an `api` struct and return a closure for one of its function fields; the caller assigns the result.

```go
// sandbox/lib/schemaitem/schemaitem.go
func CheckKeysPresenceFactory(s *api.SchemaItem) func(keys []string) bool {
	return func(keys []string) bool {
		for _, key := range keys {
			exists, err := s.Deps.Exists(dense.ValueKey(s.Prefix, s.Id, key))
			if err != nil || !exists {
				return false
			}
		}
		return true
	}
}

// build is the package's factory aggregate — the one place that must
// stay complete.
func build(d deps.Deps, items []api.Item, prefix string, id int64) api.SchemaItem {
	s := api.SchemaItem{Deps: d, Items: items, Prefix: prefix, Id: id}
	s.CheckKeysPresence = CheckKeysPresenceFactory(&s)
	// … one assignment per remaining function field
	return s
}
```

The closure captures `s`, so `s.Deps.Exists(...)` is resolved when `CheckKeysPresence` is *called*, not when the factory ran. That is what carries the injected deps into behavior the caller can hold, while `sandbox/lib/` stays unreachable from outside.

Two properties follow, and both are load-bearing:

- **One field, one factory.** A factory does nothing but return the closure; the package's constructor (`New`, or an aggregate like `build`) is the only place that has to assign them all.
- **Deps are read through the pointer, never copied into the closure.** Capturing `d` directly would freeze the dependency at construction; reading `s.Deps` keeps the struct authoritative.

---

## Adapters Fill Their Contract the Same Way

The pattern does not stop at the sandbox wall. An adapter fills `deps.Deps` with factories too — only the **carrier** changes: instead of an `api` struct holding the deps it reads, the carrier is the adapter struct, holding the configuration its closures read and declaring the contract they fill.

```go
// adapters/standard/standard.go
type StandardAdapter struct {
	Deps deps.Deps // the contract the factories assign into
	mu   sync.Mutex
	base string // the state the closures read
}

func ReadFactory(s *StandardAdapter) func(key string) ([]byte, error) {
	return func(key string) ([]byte, error) {
		value, err := os.ReadFile(s.path(key))
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%w: %s", deps.ErrKeyNotFound, key)
		}
		return value, err
	}
}

// NewWithBase is the adapter's factory aggregate: it returns the contract
// struct, never the concrete adapter type.
func NewWithBase(base string) deps.Deps {
	s := &StandardAdapter{base: base}
	s.Deps.Write = WriteFactory(s)
	s.Deps.Read = ReadFactory(s)
	// … one assignment per remaining field of deps.Deps
	return s.Deps
}
```

Binding a method into a field would work in Go, but the project does not do it: one shape for filling struct contracts means one place to look for completeness — the constructor at the bottom of the file — on both sides of the wall.

---

## Replacing One Behavior

With an interface, overriding a single method means declaring a type that embeds the original and shadows one method. With a struct, it is an assignment:

```go
myDeps := standard.New()

// Keep every other field of the adapter; disable deletes.
myDeps.Delete = func(key string) error {
	return fmt.Errorf("deletes are disabled")
}

keep := keeplib.New(myDeps)
```

This is the everyday testing path: stub one call, leave the rest of the adapter alone.

---

## No Nil Form for a Struct

`SchemaInstance` and `SchemaItem` are structs, not interfaces, so a lookup that can fail cannot return a typed `nil` to signal "not found" — a zero-value struct is never `== nil`. Every lookup that can miss returns an extra `bool` instead:

```go
// api.KeepDatabase.GetSchema and api.SchemaInstance.FindByKey
GetSchema func(name string) (api.SchemaInstance, bool)
FindByKey func(key string, keyValue any) (api.SchemaItem, bool)
```

Callers check `ok`, never `== nil`:

```go
users, ok := db.GetSchema("user")
if !ok {
	panic("schema not declared in Props")
}
```

`*api.Error`, by contrast, keeps its pointer: a nil `*Error` still reads as "no failure" everywhere `if err != nil` is used, and callers read its fields (`err.Type`, `err.Key`, `err.KeyValue`, `err.Message`) directly since `Error` carries no methods.

---

## What It Costs

The compiler no longer checks completeness. An interface implementation that misses a method fails to build; a struct contract with an unfilled field compiles fine and panics on the first call with a nil-pointer dereference.

That moves one guarantee from the compiler to the author:

- An adapter's `New` must call a factory for **every** field of `deps.Deps`.
- Every object's constructor (`database.New`, `schemainstance.New`, the `build` aggregate in `schemaitem`) must call **every** field factory of the `api` struct it builds.
- Adding a field to a contract means visiting every adapter — see [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md#add-a-dependency).

There is a second cost, specific to factories: **the `Deps` field is read-only once the struct is returned.** The closures capture the struct the factories ran over, so patching a copy has no effect on behavior — patch the `deps.Deps` value **before** calling `lib.New`.

In exchange, there is no partial-implementation ambiguity at the call site: a filled contract is a value that can be copied, patched field by field, and passed on.
