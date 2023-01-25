package gopkg_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"

	"github.com/thecodedproject/msgen/parser/gopkg"
)

func TestGenerateStructDecl(t *testing.T) {

	testCases := []struct{
		Name string
		S gopkg.DeclStruct
		ImportAliases map[string]string
		ExpectedErr error
	}{
		{
			Name: "empty struct",
			ExpectedErr: errors.New("struct name cannot be empty"),
		},
		{
			Name: "empty struct with name",
			S: gopkg.DeclStruct{
				Name: "MyObject",
			},
		},
/*
		{
			Name: "struct with fields no import aliases",
			S: gopkg.DeclStruct{
				Name: "MyObject",
				Fields: map[string]
			},
		},
*/
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			buffer := bytes.NewBuffer(nil)

			err := gopkg.GenerateStructDecl(
				buffer,
				test.S,
				test.ImportAliases,
			)

			if test.ExpectedErr != nil {
				require.Equal(t, test.ExpectedErr, err)
				return
			}

			require.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, t.Name(), buffer.Bytes())
		})
	}

}
