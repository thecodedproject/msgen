package main

import(
	"flag"
	"github.com/thecodedproject/msgen/stategen/generator"
	"log"
	"os"
)

var inputFile = flag.String("inputFile", "./state.go", "Input file containing input struct")
var inputStructName = flag.String("inputStruct", "stateImpl", "Name of input struct")
var outputInterfaceNameCmd = flag.String("outputInterface", "", "Name of output interface (Or do not generate interface if not set)")
var outputFile = flag.String("outputFile", "./state_gen.go", "Output file to write output interface to")

func main() {

	flag.Parse()

	CheckInputFileExists(*inputFile)

	var outputInterfaceName *string
	if *outputInterfaceNameCmd != "" {
		outputInterfaceName = outputInterfaceNameCmd
	}

	err := generator.Generate(
		*inputFile,
		*inputStructName,
		outputInterfaceName,
		*outputFile,
	)
	if err != nil {
		log.Fatal("stategen: ", err.Error())
	}
}

func CheckInputFileExists(inputFile string) {

	info, err := os.Stat(inputFile)
	if os.IsNotExist(err) {
		log.Fatalf("stategen: Input file '%s' does not exist", inputFile)
	} else if err != nil {
		log.Fatal("stategen: Cannot stat input file", inputFile, err.Error())
	}

	if info.IsDir() {
		log.Fatalf("stategen: Proto file '%s' is a directory", inputFile)
	}
}
