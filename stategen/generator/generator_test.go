package generator_test

import(
	"bytes"
	"flag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thecodedproject/msgen/stategen/generator"
	"io/ioutil"
	"os"
	"testing"
)

var fix = flag.Bool(
	"fix",
	false,
	"Overwrite expected output files with actual test results",
)

func TestGenerator(t *testing.T) {

	testCases := []struct{
		Name string
		InputStructName string
		OutputInterfaceName string
		InputFile string
		ExpectedOutputFile string
	}{
		{
			Name: "Empty state",
			InputStructName: "inputImpl",
			OutputInterfaceName: "OutputInterface",
			InputFile: "./test_files/test_generator_empty_state.in",
			ExpectedOutputFile: "./test_files/test_generator_empty_state.out",
		},
		{
			Name: "State struct with multiple members and New method",
			InputStructName: "stateImpl",
			OutputInterfaceName: "State",
			InputFile: "./test_files/test_generator_state_struct_with_multiple_members_and_new_method.in",
			ExpectedOutputFile: "./test_files/test_generator_state_struct_with_multiple_members_and_new_method.out",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

				buffer := bytes.NewBuffer(nil)

				err := generator.Generate(
					test.InputFile,
					test.InputStructName,
					test.OutputInterfaceName,
					buffer,
				)
				require.NoError(t, err)

				if *fix {
					outFile, err := os.Create(test.ExpectedOutputFile)
					require.NoError(t, err)
					defer outFile.Close()

					outFile.Write(buffer.Bytes())
					return
				}

				expectedBuffer, err := ioutil.ReadFile(test.ExpectedOutputFile)
				require.NoError(t, err)

				assert.Equal(t, string(expectedBuffer), buffer.String())
		})
	}
}
