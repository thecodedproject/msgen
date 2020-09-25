package proto_type_conversion

import(
	"github.com/thecodedproject/msgen/generator/files/common"
	"github.com/thecodedproject/msgen/parser"
	"github.com/thecodedproject/msgen/parser/proto_helpers"
	"io"
	"path"
)

const(
	filename = "type_conversion.go"
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

	nestedTypes := proto_helpers.NestedMessages(i)

	if len(nestedTypes) == 0 {
		return nil
	}

	_, serviceName := path.Split(serviceRootImportPath)
	outputFile := path.Join(outputDir, serviceName + "pb", filename)

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

	nestedTypes := proto_helpers.NestedMessages(i)

	if len(nestedTypes) == 0 {
		return nil
	}

	baseTemplate := common.BaseTemplate()

	header, err := baseTemplate.Parse(headerTmpl)
	if err != nil {
		return err
	}

	_, serviceName := path.Split(serviceRootImportPath)

	return header.Execute(writer, struct{
		Package string
		Imports []string
		TypesPackage string
		NestedTypes []parser.Message
	}{
		Package: serviceName + "pb",
		Imports: []string{
			"\"" + serviceRootImportPath + "\"",
		},
		TypesPackage: serviceName,
		NestedTypes: nestedTypes,
	})
}

var headerTmpl = `
{{- $typesPackage := .TypesPackage -}}
package {{.Package}}

// Generated by msgen. DO NOT EDIT.

import(
{{- range .Imports}}
	{{.}}
{{- end}}
)

{{- range .NestedTypes}}
{{- $typeName := print (ToCamel .Name)}}

func {{$typeName}}FromProto(v *{{$typeName}}) *{{$typesPackage}}.{{$typeName}} {

	if v == nil {
		return nil
	}

	return &{{$typesPackage}}.{{$typeName}}{
{{- range .Fields}}
		{{ToCamel .Name}}: v.{{ToCamel .Name}},
{{- end}}
	}
}

func {{$typeName}}ToProto(v *{{$typesPackage}}.{{$typeName}}) *{{$typeName}} {

	if v == nil {
		return nil
	}

	return &{{$typeName}}{
{{- range .Fields}}
		{{ToCamel .Name}}: v.{{ToCamel .Name}},
{{- end}}
	}
}
{{- end}}

`
