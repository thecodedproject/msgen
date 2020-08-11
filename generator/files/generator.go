package files

import(
	"github.com/thecodedproject/msgen/generator/files/client_logical"
	"github.com/thecodedproject/msgen/parser"
	//"path"
)

func GenerateFiles(
	serviceRootImportPath string,
	i parser.ProtoInterface,
	outputDir string,
) error {

	//... other generators...

	err := client_logical.Generate(
		serviceRootImportPath,
		i,
		outputDir,
	)
	if err != nil {
		return err
	}

	return nil
}
