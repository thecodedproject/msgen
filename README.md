
## Examples

An example microservice is located in the `examples/pinger` folder. This includes only the `.proto` definition file and a Golang `generate.go` file.

To generate the microservice, run:

```
cd examples/pinger
go generate generate.go
```

Tests are generated for all client interface functions; to run these, do:

```
go test ./...
```

(Note the tests fail with a `TODO: Implement test...` message)

Clean any of the example folders (i.e. remove all generated files) with:

```
git clean -xfd
```
