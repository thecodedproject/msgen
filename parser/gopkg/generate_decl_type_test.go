package gopkg_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"

	"github.com/thecodedproject/msgen/parser/gopkg"
)

func TestGenerateDeclType(t *testing.T) {

	testCases := []struct{
		Name string
		T gopkg.DeclType
		ImportAliases map[string]string
		ExpectedErr error
	}{
		{
			Name: "empty name return error",
			ExpectedErr: errors.New("type decl name cannot be empty"),
		},
		{
			Name: "empty type returns error",
			T: gopkg.DeclType{
				Name: "SomeTypeName",
			},
			ExpectedErr: errors.New("type decl type cannot be nil"),
		},
		{
			Name: "TypeStruct with no fields",
			T: gopkg.DeclType{
				Name: "MyStrct",
				Type: gopkg.TypeStruct{},
			},
		},
		{
			Name: "TypeStruct with fields no import aliases",
			T: gopkg.DeclType{
				Name: "MyStrct",
				Type: gopkg.TypeStruct{
					Fields: []gopkg.DeclVar{
						{
							Name: "SomeValue",
							Type: gopkg.TypeInt{},
						},
						{
							Name: "SomeOtherValue",
							Type: gopkg.TypeFloat64{},
						},
					},
				},
			},
		},
		{
			Name: "TypeStruct with fields and import aliases",
			T: gopkg.DeclType{
				Name: "MyStrct",
				Type: gopkg.TypeStruct{
					Fields: []gopkg.DeclVar{
						{
							Name: "MapOfThings",
							Type: gopkg.TypeMap{
								KeyType: gopkg.TypeString{},
								ValueType: gopkg.TypeUnknownNamed{
									Name: "CustomContainer",
									Import: "some/pkgpath",
								},
							},
						},
						{
							Name: "OtherThing",
							Type: gopkg.TypeUnknownNamed{
								Name: "ThingType",
								Import: "some/otherimport",
							},
						},
					},
				},
			},
			ImportAliases: map[string]string{
				"some/otherimport": "otherimport",
				"some/pkgpath": "pkgpath",
			},
		},
		{
			Name: "TypeMap with import aliases",
			T: gopkg.DeclType{
				Name: "MyCustomMap",
				Type: gopkg.TypeMap{
					KeyType: gopkg.TypeString{},
					ValueType: gopkg.TypeUnknownNamed{
						Name: "CustomContainer",
						Import: "some/pkgpath",
					},
				},
			},
			ImportAliases: map[string]string{
				"some/pkgpath": "pkgpath",
			},
		},
		{
			Name: "TypeInterface with no fields",
			T: gopkg.DeclType{
				Name: "MyServiceApi",
				Type: gopkg.TypeInterface{},
			},
		},
		{
			Name: "TypeInterface with no fields",
			T: gopkg.DeclType{
				Name: "MyServiceApi",
				Type: gopkg.TypeInterface{},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			buffer := bytes.NewBuffer(nil)

			err := gopkg.GenerateDeclType(
				buffer,
				test.T,
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
