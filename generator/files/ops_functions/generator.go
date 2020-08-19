package ops_functions

import(
	"github.com/iancoleman/strcase"
	"github.com/thecodedproject/msgen/generator/files/common"
	"github.com/thecodedproject/msgen/generator/files/proto_helpers"
	"github.com/thecodedproject/msgen/parser"
	"io"
	"path"
)

const(
	opsPath = "ops"
)

type Method struct {
	Name string
	Args []parser.Field
	ReturnArgs []parser.Field
}

func Generate(
	serviceRootImportPath string,
	i parser.ProtoInterface,
	outputDir string,
) error {


	for _, method := range i.Methods {

		outputFile := path.Join(
			outputDir,
			opsPath,
			strcase.ToSnake(method.Name) + ".go",
		)

		writer, err := common.CreatePathAndOpen(outputFile)
		if err != nil {
			return err
		}

		err = GenerateBufferForMethod(
			serviceRootImportPath,
			i,
			writer,
			method.Name,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateBufferForMethod(
	serviceRootImportPath string,
	i parser.ProtoInterface,
	writer io.Writer,
	methodName string,
) error {

	baseTemplate := common.BaseTemplate()

	header, err := baseTemplate.Parse(headerTmpl)
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

	methodTemplate, err := baseTemplate.Parse(methodTmpl)
	if err != nil {
		return err
	}

	err = methodTemplate.Execute(writer, methodParams)
	if err != nil {
		return err
	}

	return nil
}

var headerTmpl = `package {{.Package}}

import(
{{- range .Imports}}
	{{.}}
{{- end}}
)

`

var methodTmpl = `func {{ToCamel .Name}}(
	ctx context.Context,
	b Backends,
{{- range .Args}}
	{{ToLowerCamel .Name}} {{.Type}},
{{- end}}
) {{FuncRetValsWithError .ReturnArgs}} {

	{{FuncDefaultReturn_WithError .ReturnArgs}}
}

`

