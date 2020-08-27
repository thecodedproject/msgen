package state

import(
	"github.com/thecodedproject/msgen/generator/files/common"
	"io"
	"path"
	stategen "github.com/thecodedproject/msgen/stategen/generator"
)

const(
	stateFolder = "state"

	// Parameters for stategen generator
	stategenInputFile = "state.go"
	stategenInputStructName = "stateImpl"
	stategenOutputInterfaceName = "State"
	stategenOutputFile = "state_gen.go"
)

func Generate(
	serviceRootImportPath string,
	outputDir string,
) error {

	stateFile := path.Join(outputDir, stateFolder, stategenInputFile)

	stateFileExists, err := common.FileExists(stateFile)
	if err != nil {
		return err
	}

	if !stateFileExists {
		writer, err := common.CreatePathAndOpen(stateFile)
		if err != nil {
			return err
		}

		err = GenerateBuffer(
			serviceRootImportPath,
			writer,
		)
		if err != nil {
			return err
		}
	}

	stateGenFile := path.Join(outputDir, stateFolder, stategenOutputFile)
	stateGenFileExists, err := common.FileExists(stateGenFile)

	if stateGenFileExists {
		return nil
	}

	return stategen.Generate(
		stateFile,
		stategenInputStructName,
		stategenOutputInterfaceName,
		stateGenFile,
	)
}

func GenerateBuffer(
	serviceRootImportPath string,
	writer io.Writer,
) error {

	baseTemplate := common.BaseTemplate()

	header, err := baseTemplate.Parse(headerTmpl)
	if err != nil {
		return err
	}

	return header.Execute(writer, struct{
		Package string
		InputFile string
		InputStruct string
		OutputInterface string
		OutputFile string
	}{
		Package: stateFolder,
		InputFile: stategenInputFile,
		InputStruct: stategenInputStructName,
		OutputInterface: stategenOutputInterfaceName,
		OutputFile: stategenOutputFile,
	})
}

var headerTmpl = `package {{.Package}}

//go:generate stategen --inputFile="{{.InputFile}}" --inputStruct="{{.InputStruct}}" --outputInterface="{{.OutputInterface}}" --outputFile="{{.OutputFile}}"

import(
)

type stateImpl struct {
}

func New() *stateImpl {

	return &stateImpl{
	}
}

`

