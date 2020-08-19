package parser_test

import(
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thecodedproject/msgen/parser"
)

func TestParse(t *testing.T) {

	testCases := []struct{
		Name string
		ProtoFile string
		Expected parser.ProtoInterface
	}{
		{
			Name: "some",
			ProtoFile: "./example_proto/TestParse_some.proto",
			Expected: parser.ProtoInterface{
				ProtoPackage: "examplepb",
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
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			i, err := parser.Parse(test.ProtoFile)
			require.NoError(t, err)

			assert.Equal(t, test.Expected, i)
		})
	}
}
