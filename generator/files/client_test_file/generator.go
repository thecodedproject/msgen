package client_test_file

import(
	"github.com/thecodedproject/msgen/generator/files/common"
	"github.com/thecodedproject/msgen/parser"
	"github.com/thecodedproject/msgen/parser/proto_helpers"
	"io"
	"path"
)

const(
	relativePath = "client/client_test.go"
)

type Method struct {
	Name string
	ServiceName string
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

	header, err := baseTemplate.Parse(testHeaderTmpl)
	if err != nil {
		return err
	}

	err = header.Execute(writer, struct{
		Package string
		Imports []string
		ServiceName string
	}{
		Package: "client_test",
		Imports: []string{
			"\"context\"",
			"\"github.com/stretchr/testify/suite\"",
			"\"google.golang.org/grpc\"",
			"\"google.golang.org/grpc/connectivity\"",
			"\"" + serviceRootImportPath + "\"",
			"\"" + serviceRootImportPath + "/" + i.ProtoPackage + "\"",
			"\"" + serviceRootImportPath + "/rpc_server\"",
			"\"" + serviceRootImportPath + "/state\"",
			"logical_client \"" + serviceRootImportPath + "/client/logical\"",
			"grpc_client \"" + serviceRootImportPath + "/client/grpc\"",
			"\"testing\"",
			"\"time\"",
			"\"net\"",
			"\"log\"",
		},
		ServiceName: serviceName,
	})
	if err != nil {
		return err
	}


	for _, method := range i.Methods {

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
			ServiceName: serviceName,
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
	}

	return nil
}

var testHeaderTmpl = `package {{.Package}}

// Generated by msgen. DO NOT EDIT.

import(
{{- range .Imports}}
	{{.}}
{{- end}}
)

func setupGRPCServer(ts *TestGRPCSuite, s state.State) (string) {

	listener, err := net.Listen("tcp", "localhost:0")
	ts.Require().NoError(err)

	grpcSrv := grpc.NewServer()
	ts.T().Cleanup(grpcSrv.GracefulStop)

	{{ToLowerCamel .ServiceName}}Srv := rpc_server.New(s)
	{{ToLowerCamel .ServiceName}}pb.Register{{ToCamel .ServiceName}}Server(grpcSrv, {{ToLowerCamel .ServiceName}}Srv)

	go func() {
		err := grpcSrv.Serve(listener)
		ts.Require().NoError(err)
	}()

	return listener.Addr().String()
}

func setupGRPCClient(ts *TestGRPCSuite, s state.State) {{ToLower .ServiceName}}.Client {

	serverAddr := setupGRPCServer(ts, s)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	ts.Require().NoError(err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for {
		if conn.GetState() == connectivity.Ready {
			break
		}

		if !conn.WaitForStateChange(ctx, conn.GetState()) {
			log.Fatal("grpc timeout whilst connecting")
		}
	}

	client := grpc_client.NewForTesting(ts.T(), conn)
	return client
}


type clientSuite struct {
	suite.Suite

	createClient func(state.State) {{ToLower .ServiceName}}.Client
}

func TestLogical(t *testing.T) {
	suite.Run(t, new(TestLogicalSuite))
}

func TestGRPC(t *testing.T) {
	suite.Run(t, new(TestGRPCSuite))
}

type TestLogicalSuite struct {
	clientSuite
}

func (ts *TestLogicalSuite) SetupTest() {
	ts.createClient = func(s state.State) {{ToLower .ServiceName}}.Client {
		return logical_client.New(s)
	}
}

type TestGRPCSuite struct {
	clientSuite
}

func (ts *TestGRPCSuite) SetupTest() {
	ts.createClient = func(s state.State) {{ToLower .ServiceName}}.Client {
		return setupGRPCClient(ts, s)
	}
}

`

var testMethodTmpl = `
{{- $serviceName := .ServiceName -}}
func (ts *clientSuite) Test{{ToCamel .Name}}() {

	testCases := []struct{
		Name string
	}{
		{
			Name: "some_test",
		},
	}

	for _, test := range testCases {
		ts.T().Run(test.Name, func(t *testing.T) {

			c := ts.createClient(
				state.NewStateForTesting(
					t,
				),
			)

			ctx := context.Background()
{{- range .Args}}
{{- if .IsNestedMessage}}
			var {{ToLowerCamel .Name}} {{ToLower $serviceName}}.{{.Type}}
{{- else}}
			var {{ToLowerCamel .Name}} {{.Type}}
{{- end}}
{{- end}}
			var err error
			{{range .ReturnArgs}}_, {{end}}err = c.{{ToCamel .Name}}(
				ctx,
{{- range .Args}}
				{{if .IsNestedMessage}}&{{end}}{{ToLowerCamel .Name}},
{{- end}}
			)
			ts.Require().NoError(err)

			ts.Assert().Fail("TODO: Implement test for client.{{ToCamel .Name}}")
		})
	}
}

`

