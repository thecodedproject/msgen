package gopkg_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thecodedproject/msgen/parser/gopkg"
)

func TestTypeDefaultInit(t *testing.T) {

	testCases := []struct{
		Def gopkg.Type
		ImportAliases map[string]string
		Expected string
		ExpectedErr error
	}{
		{
			Def: gopkg.TypeByte{},
			Expected: "0",
		},
		{
			Def: gopkg.TypeError{},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeFloat32{},
			Expected: "0",
		},
		{
			Def: gopkg.TypeFloat64{},
			Expected: "0",
		},
		{
			Def: gopkg.TypeInt{},
			Expected: "0",
		},
		{
			Def: gopkg.TypeInt32{},
			Expected: "0",
		},
		{
			Def: gopkg.TypeInt64{},
			Expected: "0",
		},
		{
			Def: gopkg.TypeString{},
			Expected: "\"\"",
		},
		{
			Def: gopkg.TypeInterface{
				Name: "MyInterfaceTypeNoAlias",
				Import: "some_import",
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeInterface{
				Name: "MyInterfaceType",
				Import: "some/other/import",
			},
			ImportAliases: map[string]string{
				"some/other/import": "some_other_alias",
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeStruct{
				Name: "MyStructTypeNoAlias",
				Import: "some_import",
			},
			Expected: "MyStructTypeNoAlias{}",
		},
		{
			Def: gopkg.TypeStruct{
				Name: "MyStructTypeWithAlias",
				Import: "some_import",
			},
			ImportAliases: map[string]string{
				"some_import": "some_alias",
			},
			Expected: "some_alias.MyStructTypeWithAlias{}",
		},
		{
			Def: gopkg.TypeArray{
				ValueType: gopkg.TypeUnknownNamed{
					Name: "SomeType",
					Import: "some/import",
				},
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypePointer{
				ValueType: gopkg.TypeUnknownNamed{
					Name: "float32",
				},
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeMap{
				KeyType: gopkg.TypeUnknownNamed{
					Name: "MyType",
					Import: "my/import/path",
				},
				ValueType: gopkg.TypePointer{
					ValueType: gopkg.TypeArray{
						ValueType: gopkg.TypeUnknownNamed{
							Name: "MyOtherType",
							Import: "other/import",
						},
					},
				},
			},
			Expected: "nil",
		},
		{
			Def: gopkg.TypeUnknownNamed{
				Name: "MyType",
				Import: "my/import/path",
			},
			ExpectedErr: errors.New("not implemented"),
		},
	}

	for _, test := range testCases {
		t.Run(test.Def.FullType(nil), func(t *testing.T) {
			actual, err := test.Def.DefaultInit(test.ImportAliases)

			if test.ExpectedErr != nil {
				require.Equal(t, test.ExpectedErr, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.Expected, actual)
		})
	}
}

func TestTypeFullType(t *testing.T) {

	testCases := []struct{
		Def gopkg.Type
		ImportAliases map[string]string
		Expected string
	}{
		{
			Def: gopkg.TypeByte{},
			Expected: "byte",
		},
		{
			Def: gopkg.TypeError{},
			Expected: "error",
		},
		{
			Def: gopkg.TypeFloat32{},
			Expected: "float32",
		},
		{
			Def: gopkg.TypeFloat64{},
			Expected: "float64",
		},
		{
			Def: gopkg.TypeInt{},
			Expected: "int",
		},
		{
			Def: gopkg.TypeInt32{},
			Expected: "int32",
		},
		{
			Def: gopkg.TypeInt64{},
			Expected: "int64",
		},
		{
			Def: gopkg.TypeString{},
			Expected: "string",
		},
		{
			Def: gopkg.TypeInterface{
				Name: "MyInterfaceType",
				Import: "some_import",
			},
			Expected: "MyInterfaceType",
		},
		{
			Def: gopkg.TypeInterface{
				Name: "MyInterfaceType",
				Import: "some/other/import",
			},
			ImportAliases: map[string]string{
				"some/other/import": "some_other_alias",
			},
			Expected: "some_other_alias.MyInterfaceType",
		},
		{
			Def: gopkg.TypeStruct{
				Name: "MyType",
				Import: "some_import",
			},
			Expected: "MyType",
		},
		{
			Def: gopkg.TypeStruct{
				Name: "MyType",
				Import: "some_import",
			},
			ImportAliases: map[string]string{
				"some_import": "some_alias",
			},
			Expected: "some_alias.MyType",
		},
		{
			Def: gopkg.TypeArray{
				ValueType: gopkg.TypeUnknownNamed{
					Name: "SomeType",
					Import: "some/import",
				},
			},
			ImportAliases: map[string]string{
				"some/import": "some_alias",
			},
			Expected: "[]some_alias.SomeType",
		},
		{
			Def: gopkg.TypePointer{
				ValueType: gopkg.TypeUnknownNamed{
					Name: "float32",
				},
			},
			Expected: "*float32",
		},
		{
			Def: gopkg.TypeMap{
				KeyType: gopkg.TypeUnknownNamed{
					Name: "MyType",
					Import: "my/import/path",
				},
				ValueType: gopkg.TypePointer{
					ValueType: gopkg.TypeArray{
						ValueType: gopkg.TypeUnknownNamed{
							Name: "MyOtherType",
							Import: "other/import",
						},
					},
				},
			},
			ImportAliases: map[string]string{
				"my/import/path": "path_alias",
				"other/import": "other_alias",
			},
			Expected: "map[path_alias.MyType]*[]other_alias.MyOtherType",
		},
	}

	for _, test := range testCases {
		t.Run(test.Expected, func(t *testing.T) {
			require.Equal(t, test.Expected, test.Def.FullType(test.ImportAliases))
		})
	}
}
