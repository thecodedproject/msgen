package ops_backends

import(
	"github.com/thecodedproject/msgen/generator/files/common"
	"github.com/thecodedproject/msgen/generator/files/proto_helpers"
	"github.com/thecodedproject/msgen/parser"
	"io"
)

func Generate(
	outputDir string,
) error {

	panic("not implemented")
	return nil
}

func GenerateBuffer(
	serviceRootImportPath string,
	writer io.Writer,
) error {

	baseTemplate := common.BaseTemplate()

	header, err := baseTemplate.Parse(testHeaderTmpl)
	if err != nil {
		return err
	}

	return header.Execute(writer, struct{})
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

