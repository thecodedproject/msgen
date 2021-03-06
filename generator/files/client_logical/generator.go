package client_logical

import(
	"github.com/thecodedproject/msgen/generator/files/common"
	"github.com/thecodedproject/msgen/parser"
	"github.com/thecodedproject/msgen/parser/proto_helpers"
	"io"
	"path"
)

const(
	relativePath = "client/logical/client.go"
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

	outputFile := path.Join(outputDir, relativePath)

	writer, err := common.CreatePathAndOpen(outputFile)
	if err != nil {
		return err
	}

	return GenerateBuffer(
		serviceRootImportPath,
		i,
		writer,
	)
}

func GenerateBuffer(
	serviceRootImportPath string,
	i parser.ProtoInterface,
	writer io.Writer,
) error {

	serviceName := common.ServiceNameFromRootImportPath(serviceRootImportPath)

	baseTemplate := common.BaseTemplate(serviceName)

	header, err := baseTemplate.Parse(logicalHeaderTmpl)
	if err != nil {
		return err
	}

	err = header.Execute(writer, struct{
		Package string
		Imports []string
	}{
		Package: "logical",
		Imports: common.SortedImportsWithNestedTypesImport(
			serviceRootImportPath,
			i,
			"context",
			serviceRootImportPath + "/ops",
			serviceRootImportPath + "/state",
		),
	})
	if err != nil {
		return err
	}

	for _, method := range i.Methods {

		methodTemplate, err := baseTemplate.Parse(logicalMethodTmpl)
		if err != nil {
			return err
		}

		args, err := proto_helpers.MethodRequestFields(i, method.Name)
		if err != nil {
			return err
		}

		returnArgs, err := proto_helpers.MethodResponseFields(i, method.Name)
		if err != nil {
			return err
		}

		methodParams := Method{
			Name: method.Name,
			Args: args,
			ReturnArgs: returnArgs,
		}

		methodParams.ReturnArgs = append(methodParams.ReturnArgs, parser.Field{
			Name: "err",
			Type: "error",
		})

		err = methodTemplate.Execute(writer, methodParams)
		if err != nil {
			return err
		}
	}

	return nil

}

var logicalHeaderTmpl = `package {{.Package}}

// Generated by msgen. DO NOT EDIT.

import(
{{- range .Imports}}
	"{{.}}"
{{- end}}
)

type client struct {
	s state.State
}

func New(s state.State) *client {

	return &client{
		s: s,
	}
}

`



var logicalMethodTmpl = `func (c *client) {{.Name}}{{SplitFuncArgsWithCtx .Args}} {{FuncRetVals .ReturnArgs}} {

	return ops.{{.Name}}(
		ctx,
		c.s,
{{- range .Args}}
		{{ToLowerCamel .Name}},
{{- end}}
	)
}

`

