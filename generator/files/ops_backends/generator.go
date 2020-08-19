package ops_backends

import(
	"github.com/thecodedproject/msgen/generator/files/common"
	"io"
	"path"
)

const(
	relativePath = "ops/backends.go"
)

func Generate(
	serviceRootImportPath string,
	outputDir string,
) error {

	outputFile := path.Join(outputDir, relativePath)

	writer, err := common.CreatePathAndOpen(outputFile)
	if err != nil {
		return err
	}

	return GenerateBuffer(
		serviceRootImportPath,
		writer,
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
	}{
		Package: "ops",
	})
}

var headerTmpl = `package {{.Package}}

//go:generate genbackendsimpl

import(
)

type Backends interface {
}

`

