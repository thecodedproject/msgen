package main

import(
	"flag"
	"github.com/thecodedproject/msgen/generator"
	"github.com/thecodedproject/msgen/parser"
	"log"
)

var protoPath = flag.String("proto", "", "path to proto file")
var msDir = flag.String("msdir", ".", "dir to put microservice in")
var importBasePath = flag.String("import", "", "base import for microserice")

func main() {

	if protoPath == nil {
		log.Fatal("Must define proto path")
	}
	if importBasePath == nil {
		log.Fatal("Must specify base import for microserice")
	}

	iProto, err := parser.Parse(*protoPath)
	if err != nil {
		log.Fatal("Parse error:", err.Error())
	}

	err := generate.Generate(
		*importBasePath,
		iProto,
		*msDir,
	)
	if err != nil {
		log.Fatal("Generate error:", err.Error())
	}
}
