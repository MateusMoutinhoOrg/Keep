###  Dependency Injection 

These Documentation covers how to manage Dependency Injection system used by the lib 
The idea is that, the user can  choose between his **own backend implementation**, a **predefined backend implementation**.

### Interface
The required Backend is defined in [deps.go](/pkg/deps/deps.go)  its a interface covering all the required functions need to operate the lib.\

### Implementations
in the [adapters](/adapters/) is where the lib stores its own opinated BackEnds. The default implementation is located in [adapters/Standard/](/adapters/standard/) and is the implementation that will work in most cenarios.
