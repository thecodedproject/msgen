package users

//go:generate go run ../../main.go --import github.com/thecodedproject/msgen/examples/users

//go:generate protoc --go_out=plugins=grpc:. userspb/users.proto
