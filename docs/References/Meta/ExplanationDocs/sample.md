# Dependency Mechanics

## Description
Explains how the library's dependency injection works: adapters implement the `deps.Deps` interface, and the library calls those methods without knowing which adapter provided them.

---

## The Deps Contract

`deps.Deps` is an interface. Every adapter implements each method, and the library reaches dependencies only through it:

```go
type Deps interface {
    ExampleDepFunctionA() int // each method is one injectable behavior
}
```

Because the contract is an interface, the compiler rejects an adapter that forgets a method — a partial backend never reaches a consumer.

---

## Overriding a Dependency

An adapter's output is just a value satisfying the interface, so a single behavior can be swapped by **embedding** it and shadowing that one method:

```go
// Inherit every method of the embedded adapter...
type patched struct {
    deps.Deps
}

// ...and shadow only the behavior you need to change.
func (patched) ExampleDepFunctionA() int {
    return 404
}

l := lib.New(patched{Deps: standard.New(3)}) // the library is unaware of the change
```
