package parser_test

import(
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thecodedproject/msgen/parser"
	"github.com/thecodedproject/msgen/parser/methods"
)

func TestParse(t *testing.T) {

	testCases := []struct{
		Name string
		ProtoFile string
		Expected parser.ProtoInterface
	}{
		{
			Name: "Two different methods",
			ProtoFile: "./example_proto/TestParse_two_different_methods.proto",
			Expected: parser.ProtoInterface{
				ProtoPackage: "examplepb",
				ServiceName: "Some",
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
		{
			Name: "Nested messages",
			ProtoFile: "./example_proto/TestParse_nested_messages.proto",
			Expected: parser.ProtoInterface{
				ProtoPackage: "nestedpb",
				ServiceName: "NestedService",
				Methods: []parser.Method{
					{
						Name: "Ping",
						RequestMessage: "PingRequest",
						ResponseMessage: "PingResponse",
					},
				},
				Messages: []parser.Message{
					{
						Name: "PingRequest",
						Fields: []parser.Field{
							{
								Name: "some_nested_value",
								Type: "NestedVal",
								IsNestedMessage: true,
							},
						},
					},
					{
						Name: "PingResponse",
						Fields: []parser.Field{
							{
								Name: "some_other_value",
								Type: "OtherNestedVal",
								IsNestedMessage: true,
							},
						},
					},
					{
						Name: "NestedVal",
						Fields: []parser.Field{
							{
								Name: "some_value",
								Type: "int64",
							},
						},
					},
					{
						Name: "OtherNestedVal",
						Fields: []parser.Field{
							{
								Name: "some_string",
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

			i, err := methods.Parse(test.ProtoFile)
			require.NoError(t, err)

			assert.Equal(t, test.Expected, i)
		})
	}
}

func TestGetPackgeContents(t *testing.T) {


	packageDir := "example_proto/custom_conversion_1"

	_ = methods.AddCustomConversionFields(
		packageDir,
		[]string{"IntAsString"},
	)

	assert.Fail(t, "TODO: implement...")

}
