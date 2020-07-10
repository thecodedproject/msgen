package files_test

import(
	"bytes"
	"flag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thecodedproject/msgen/generator/files/client_logical"
	"github.com/thecodedproject/msgen/generator/files/client_grpc"
	"github.com/thecodedproject/msgen/parser"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

var fix = flag.Bool(
	"fix",
	false,
	"Overwrite expected output files with actual test results",
)

func TestGenerate(t *testing.T) {

	generatorsToTest := []struct{
		Name string
		Function func(string, parser.ProtoInterface, io.Writer) error
		ExpectedFileSuffix string
	}{
		{
			Name: "Client Logical",
			Function: client_logical.GenerateBuffer,
			ExpectedFileSuffix: "_client_logical.txt",
		},
		{
			Name: "Client GRPC",
			Function: client_grpc.GenerateBuffer,
			ExpectedFileSuffix: "_client_grpc.txt",
		},
	}

	testCases := []struct{
		Name string
		ServiceRootImportPath string
		ProtoInterface parser.ProtoInterface
		ExpectedFilePrefix string
	}{
		{
			Name: "Methods using only built in types",
			ServiceRootImportPath: "some/service",
			ProtoInterface: parser.ProtoInterface{
				Methods: []parser.Method{
					{
						Name: "Ping",
						RequestMessage: "PingRequest",
						ResponseMessage: "PingResponse",
					},
					{
						Name: "Pong",
						RequestMessage: "PongRequest",
						ResponseMessage: "PongResponse",
					},
				},
				Messages: []parser.Message{
					{
						Name: "PingRequest",
						Fields: []parser.Field{
							{
								Name: "int64_value",
								Type: "int64",
							},
							{
								Name: "string_value",
								Type: "string",
							},
						},
					},
					{
						Name: "PingResponse",
					},
					{
						Name: "PongRequest",
					},
					{
						Name: "PongResponse",
						Fields: []parser.Field{
							{
								Name: "int64_value",
								Type: "int64",
							},
							{
								Name: "string_value",
								Type: "string",
							},
						},
					},
				},
			},
			ExpectedFilePrefix: "./test_files/TestGenerate_only_built_in_types",
		},
	}

	for _, generator := range generatorsToTest {

		for _, test := range testCases {
			t.Run(test.Name + " " + generator.Name, func(t *testing.T) {

				buffer := bytes.NewBuffer(nil)

				err := generator.Function(
					test.ServiceRootImportPath,
					test.ProtoInterface,
					buffer,
				)
				require.NoError(t, err)

				expectedFilePath := test.ExpectedFilePrefix + generator.ExpectedFileSuffix

				if *fix {
					outFile, err := os.Create(expectedFilePath)
					require.NoError(t, err)
					defer outFile.Close()

					outFile.Write(buffer.Bytes())
					return
				}

				expectedBuffer, err := ioutil.ReadFile(expectedFilePath)
				require.NoError(t, err)

				assert.Equal(t, string(expectedBuffer), buffer.String())
			})
		}
	}
}
