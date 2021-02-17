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
		DontGeneratedOutputInterface bool
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
		{
			Name: "State struct with chan members",
			InputStructName: "stateImpl",
			OutputInterfaceName: "State",
			InputFile: "./test_files/test_generator_state_struct_with_chan_members.in",
			ExpectedOutputFile: "./test_files/test_generator_state_struct_with_chan_members.out",
		},
		{
			Name: "State struct with slice members",
			InputStructName: "stateImpl",
			OutputInterfaceName: "State",
			InputFile: "./test_files/test_generator_state_struct_with_slice_members.in",
			ExpectedOutputFile: "./test_files/test_generator_state_struct_with_slice_members.out",
		},
		{
			Name: "State struct with pointer members",
			InputStructName: "stateImpl",
			OutputInterfaceName: "State",
			InputFile: "./test_files/test_generator_state_struct_with_pointer_members.in",
			ExpectedOutputFile: "./test_files/test_generator_state_struct_with_pointer_members.out",
		},
		{
			Name: "Input struct with aliased imports",
			InputStructName: "stateImpl",
			OutputInterfaceName: "State",
			InputFile: "./test_files/test_generator_input_struct_with_aliased_imports.in",
			ExpectedOutputFile: "./test_files/test_generator_input_struct_with_aliased_imports.out",
		},
		{
			Name: "Input struct with no imports",
			InputStructName: "NoImportsStruct",
			OutputInterfaceName: "NoImportsInterface",
			InputFile: "./test_files/test_generator_input_struct_with_no_imports.in",
			ExpectedOutputFile: "./test_files/test_generator_input_struct_with_no_imports.out",
		},
		{
			Name: "NotGeneratingOutputInterface",
			InputStructName: "NoInterfaceStruct",
			DontGeneratedOutputInterface: true,
			InputFile: "./test_files/test_generator_not_generating_output_interface.in",
			ExpectedOutputFile: "./test_files/test_generator_not_generating_output_interface.out",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

				var buffer bytes.Buffer

				var outputInterfaceName *string
				if !test.DontGeneratedOutputInterface {
					outputInterfaceName = &test.OutputInterfaceName
				}

				err := generator.GenerateBuffer(
					test.InputFile,
					test.InputStructName,
					outputInterfaceName,
					&buffer,
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
