package ops_functions

import(
	"github.com/thecodedproject/msgen/generator/files"
	"github.com/thecodedproject/msgen/generator/files/proto_helpers"
	"github.com/thecodedproject/msgen/parser"
	"io"
)

type Method struct {
	Name string
	Args []parser.Field
	ReturnArgs []parser.Field
}

func Generate(
	i parser.ProtoInterface,
	outputDir string,
) error {

	panic("not implemented")
	return nil
}

func GenerateBufferForMethod(
	serviceRootImportPath string,
	i parser.ProtoInterface,
	writer io.Writer,
	methodName string,
) error {

	baseTemplate := files.BaseTemplate()

	header, err := baseTemplate.Parse(testHeaderTmpl)
	if err != nil {
		return err
	}

	err = header.Execute(writer, struct{
		Package string
		Imports []string
		ServiceName string
	}{
		Package: "ops",
		Imports: []string{
			"\"context\"",
		},
		ServiceName: "SomeService",
	})
	if err != nil {
		return err
	}

	args, err := proto_helpers.MethodRequestFields(i, methodName)
	if err != nil {
		return err
	}

	returnArgs, err := proto_helpers.MethodResponseFields(i, methodName)
	if err != nil {
		return err
	}

	methodParams := Method{
		Name: methodName,
		Args: args,
		ReturnArgs: returnArgs,
	}

	methodTemplate, err := baseTemplate.Parse(testMethodTmpl)
	if err != nil {
		return err
	}

	err = methodTemplate.Execute(writer, methodParams)
	if err != nil {
		return err
	}

	return nil
}

var testHeaderTmpl = `package {{.Package}}

import(
{{- range .Imports}}
	{{.}}
{{- end}}
)

`

var testMethodTmpl = `func {{ToCamel .Name}}(
	ctx context.Context,
	b Backends,
{{- range .Args}}
	{{ToLowerCamel .Name}} {{.Type}},
{{- end}}
) {{FuncRetValsWithError .ReturnArgs}} {

	{{FuncDefaultReturn_WithError .ReturnArgs}}
}

`

