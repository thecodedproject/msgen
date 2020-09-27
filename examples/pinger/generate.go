package pinger

//go:generate go run ../../main.go --import github.com/thecodedproject/msgen/examples/pinger

//go:generate protoc --go_out=plugins=grpc:. pingerpb/pinger.proto
