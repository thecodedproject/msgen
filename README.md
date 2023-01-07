
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

## Install `stategen`

```
go get https://github.com/thecodedproject/msgen/stategen
```

## Notes on protoc and importing custom proto message definitions

Proto definitions are extensible, with messages (and other proto gubbins) able to imported into another proto file using the `import` keyword, e.g.
```
import "google/protobuf/timestamp.proto";
```

The `protoc` compiler will search for these imports in (at least) these locations:

1. Looking in the relative path `../include` from the `protoc` executable.

  This doesn't seem to be documented anywhere, but seems reasoanable - on a typical unix system `protoc` is installed in some `usr/bin` folder with the common proto files being installed into `usr/include` (e.g. `usr/include/google/protobuf`

2. By specifying `--proto_path=PATH` argument when calling `protoc` (can be be specified multiple times), e.g:
```
> protoc --proto_path='som/path' --proto_path='some/other/path'
```

3. The currenct working directory

This has implications for finding (and automatically generting code which calls) the conversion functions for these imported proto message types:

  * The proto files themselves will define a `go_package` option e.g.
```
  option go_package = "github.com/golang/protobuf/ptypes/timestamp";
```
    which shows where the corresponding golang definitions for these types should be imported from (this needs to be an importable golang path)

  * (TODO - 20221231) Still need figure out how this is to be solved - see one of a few options:
    * (probably easiest) Specify specific go packages which contain conversion functions;
      * scan these pacakges for all the available conversion functions
      * A conversion function will have a specific signature (i.e. `XToProto(X) pb.XMessage` and `XFromProto(pb.XMessage) X`)
      * Check we have all the reuqired conversion functions
      * If we don't then raise an error
