package files_test

import(
	"testing"
	"github.com/thecodedproject/msgen/parser"
	"github.com/thecodedproject/msgen/generator/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"bytes"
	"io/ioutil"
	"flag"
	"os"
)

var fix = flag.Bool(
	"fix",
	false,
	"Overwrite expected output files with actual test results",
)


func TestGenerateClientLogicalClientBuffer(t *testing.T) {

	testCases := []struct{
		Name string
		ServiceRootImportPath string
		ProtoInterface parser.ProtoInterface
		ExpectedFilePath string
	}{
		{
			Name: "Some",
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
			ExpectedFilePath: "./test_files/TestGenerateClientLogicalClientBuffer_Some.txt",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			buffer := bytes.NewBuffer(nil)

			err := files.GenerateClientLogicalClientBuffer(
				test.ServiceRootImportPath,
				test.ProtoInterface,
				buffer,
			)
			require.NoError(t, err)

			if *fix {
				outFile, err := os.Create(test.ExpectedFilePath)
				require.NoError(t, err)
				defer outFile.Close()

				outFile.Write(buffer.Bytes())
				return
			}

			expectedBuffer, err := ioutil.ReadFile(test.ExpectedFilePath)
			require.NoError(t, err)

			assert.Equal(t, string(expectedBuffer), buffer.String())
		})
	}
}
