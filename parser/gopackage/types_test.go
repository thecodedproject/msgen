package gopackage_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thecodedproject/msgen/parser/gopackage"
)

func TestTypeFullName(t *testing.T) {

	testCases := []struct{
		Def gopackage.Type
		ImportAliases map[string]string
		Expected string
	}{
		{
			Def: gopackage.TypeNamed{
				Name: "int32",
			},
			Expected: "int32",
		},
		{
			Def: gopackage.TypeNamed{
				Name: "MyType",
				Import: "some_import",
			},
			Expected: "MyType",
		},
		{
			Def: gopackage.TypeNamed{
				Name: "MyType",
				Import: "some_import",
			},
			ImportAliases: map[string]string{
				"some_import": "some_alias",
			},
			Expected: "some_alias.MyType",
		},
		{
			Def: gopackage.TypeArray{
				ValueType: gopackage.TypeNamed{
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
			Def: gopackage.TypePointer{
				ValueType: gopackage.TypeNamed{
					Name: "float32",
				},
			},
			Expected: "*float32",
		},
		{
			Def: gopackage.TypeMap{
				KeyType: gopackage.TypeNamed{
					Name: "MyType",
					Import: "my/import/path",
				},
				ValueType: gopackage.TypePointer{
					ValueType: gopackage.TypeArray{
						ValueType: gopackage.TypeNamed{
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
			require.Equal(t, test.Expected, test.Def.FullName(test.ImportAliases))
		})
	}
}
