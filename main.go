package main

import(
	"flag"
	"github.com/thecodedproject/msgen/generator/files"
	"github.com/thecodedproject/msgen/generator/files/common"
	"github.com/thecodedproject/msgen/parser"
	"log"
	"os"
	"path"
)

var msDir = flag.String("msdir", ".", "dir to put microservice in")
var rootImportPath = flag.String("import", "", "root import for microserice")

func main() {

	flag.Parse()

	if *rootImportPath == "" {
		log.Fatal("Must specify root import for microserice")
	}

	protoPath := GetProtoPathAndCheckExists(*msDir, *rootImportPath)

	iProto, err := parser.Parse(protoPath)
	if err != nil {
		log.Fatal("Parse error:", err.Error())
	}

	err = files.Generate(
		*rootImportPath,
		iProto,
		*msDir,
	)
	if err != nil {
		log.Fatal("Generate error:", err.Error())
	}
}

func GetProtoPathAndCheckExists(msDir, rootImportPath string) string {

	serviceName := common.ServiceNameFromRootImportPath(rootImportPath)

	protoPath := path.Join(msDir, serviceName + "pb", serviceName + ".proto")

	info, err := os.Stat(protoPath)
	if os.IsNotExist(err) {
		log.Fatalf("No proto file at '%s'; create this first", protoPath)
	} else if err != nil {
		log.Fatal("Cannot stat proto file", protoPath, err.Error())
	}

	if info.IsDir() {
		log.Fatalf("Proto file '%s' is a directory", protoPath)
	}

	return protoPath
}
