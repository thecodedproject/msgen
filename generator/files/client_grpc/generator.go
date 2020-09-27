package client_grpc

import(
	"github.com/thecodedproject/msgen/generator/files/common"
	"github.com/thecodedproject/msgen/parser"
	"github.com/thecodedproject/msgen/parser/proto_helpers"
	"io"
	"path"
)

const(
	relativePath = "client/grpc/client.go"
)

type Method struct {
	Name string
	Args []parser.Field
	ReturnArgs []parser.Field
	ProtoPackage string
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

	baseTemplate := common.BaseTemplate()

	header, err := baseTemplate.Parse(grpcHeaderTmpl)
	if err != nil {
		return err
	}

	serviceName := common.ServiceNameFromRootImportPath(serviceRootImportPath)

	err = header.Execute(writer, struct{
		Package string
		Imports []string
		ServiceName string
		ProtoPackage string
		ProtoServiceName string
	}{
		Package: "grpc",
		Imports: common.SortedImportsWithNestedTypesImport(
			serviceRootImportPath,
			i,
			"context",
			"errors",
			"flag",
			"google.golang.org/grpc",
			"google.golang.org/grpc/connectivity",
			serviceRootImportPath + "/" + i.ProtoPackage,
			"testing",
			"time",
		),
		ServiceName: serviceName,
		ProtoPackage: i.ProtoPackage,
		ProtoServiceName: i.ServiceName,
	})
	if err != nil {
		return err
	}


	for _, method := range i.Methods {

		args, err := proto_helpers.MethodRequestFieldsWithImportOnNestedFields(
			i,
			method.Name,
			serviceName,
		)
		if err != nil {
			return err
		}

		returnArgs, err := proto_helpers.MethodResponseFieldsWithImportOnNestedTypes(
			i,
			method.Name,
			serviceName,
		)
		if err != nil {
			return err
		}

		methodParams := Method{
			Name: method.Name,
			Args: args,
			ReturnArgs: returnArgs,
			ProtoPackage: i.ProtoPackage,
		}

		if len(methodParams.ReturnArgs) == 0 {

			methodEmptyReturnTemplate, err := baseTemplate.Parse(grpcMethodEmptyReturnTmpl)
			if err != nil {
				return err
			}

			err = methodEmptyReturnTemplate.Execute(writer, methodParams)
			if err != nil {
				return err
			}

		} else {

			methodTemplate, err := baseTemplate.Parse(grpcMethodTmpl)
			if err != nil {
				return err
			}

			err = methodTemplate.Execute(writer, methodParams)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

var grpcHeaderTmpl = `package {{.Package}}

// Generated by msgen. DO NOT EDIT.

import(
{{- range .Imports}}
	"{{.}}"
{{- end}}
)

var address = flag.String("{{ToLower .ServiceName}}_grpc_address", "", "host:port of {{ToLower .ServiceName}} gRPC service")

type client struct {
	rpcConn *grpc.ClientConn
	rpcClient {{.ProtoPackage}}.{{.ProtoServiceName}}Client
}

func IsGRPCEnabled() bool {
	return *address != ""
}

func New() (*client, error) {
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for {
		if conn.GetState() == connectivity.Ready {
			break
		}
		if !conn.WaitForStateChange(ctx, conn.GetState()) {
			return nil, errors.New("grpc timeout whilst connecting")
		}
	}

	return &client{
		rpcConn: conn,
		rpcClient: {{.ProtoPackage}}.New{{.ProtoServiceName}}Client(conn),
	}, nil
}

func NewForTesting(t *testing.T, conn *grpc.ClientConn) *client {
	return &client{
		rpcConn: conn,
		rpcClient: {{.ProtoPackage}}.New{{.ProtoServiceName}}Client(conn),
	}
}

`

var grpcMethodTmpl = `func (c *client) {{ToCamel .Name}}{{SplitFuncArgsWithCtx .Args}} {{NamedFuncRetValsWithError .ReturnArgs}} {

	res, err := c.rpcClient.{{ToCamel .Name}}(
		ctx,
		&{{.ProtoPackage}}.{{ToCamel .Name}}Request{
{{- range .Args}}
			{{ToCamel .Name}}: {{ToLowerCamel .Name}},
{{- end}}
		},
	)
	if err != nil {
		{{FuncDefaultReturn_Named_WithError .ReturnArgs}}
	}

	return {{range $index, $elem := .ReturnArgs}}{{if $index}}, {{end}}res.{{ToCamel .Name}}{{end}}, nil
}

`

var grpcMethodEmptyReturnTmpl = `func (c *client) {{ToCamel .Name}}{{SplitFuncArgsWithCtx .Args}} error {

	_, err := c.rpcClient.{{ToCamel .Name}}(
		ctx,
		&{{.ProtoPackage}}.{{ToCamel .Name}}Request{
{{- range .Args}}
			{{ToCamel .Name}}: {{ToLowerCamel .Name}},
{{- end}}
		},
	)
	return err
}

`
