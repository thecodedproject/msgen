package gopackage_test

import(
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thecodedproject/msgen/parser/gopackage"
	"testing"
)

func TestGetPackageContents(t *testing.T) {

	testCases := []struct{
		Name string
		PkgDir string
		PkgImportPath string
		Expected *gopackage.Contents
	}{
		{
			Name: "all_built_in_golang_types",
			PkgDir: "test_packages/all_built_in_types",
			PkgImportPath: "some/import/all_built_in_types",
			Expected: &gopackage.Contents{
				Functions: []gopackage.Function{
					{
						Name: "SomeInts",
						Import: "some/import/all_built_in_types",
						Args: []gopackage.Type{
							gopackage.TypeNamed{Name: "int"},
							gopackage.TypeNamed{Name: "int64"},
							gopackage.TypeNamed{Name: "int32"},
						},
						ReturnArgs: []gopackage.Type{
							gopackage.TypeNamed{Name: "int"},
							gopackage.TypeNamed{Name: "int64"},
							gopackage.TypeNamed{Name: "int32"},
						},
					},
					{
						Name: "SomeFloats",
						Import: "some/import/all_built_in_types",
						Args: []gopackage.Type{
							gopackage.TypeNamed{Name: "float32"},
							gopackage.TypeNamed{Name: "float64"},
						},
						ReturnArgs: []gopackage.Type{
							gopackage.TypeNamed{Name: "float32"},
							gopackage.TypeNamed{Name: "float64"},
						},
					},
					{
						Name: "SomeStrings",
						Import: "some/import/all_built_in_types",
						Args: []gopackage.Type{
							gopackage.TypeNamed{Name: "string"},
						},
						ReturnArgs: []gopackage.Type{
							gopackage.TypeNamed{Name: "string"},
						},
					},
				},
				StructTypes: []gopackage.StructDecl{
					{
						Name: "SomeStruct",
						Import: "some/import/all_built_in_types",
						Fields: map[string]gopackage.Type{
							"IA": gopackage.TypeNamed{Name: "int"},
							"IB": gopackage.TypeNamed{Name: "int32"},
							"IC": gopackage.TypeNamed{Name: "int64"},
							"FA": gopackage.TypeNamed{Name: "float32"},
							"FB": gopackage.TypeNamed{Name: "float64"},
							"S": gopackage.TypeNamed{Name: "string"},
						},
					},
				},
			},
		},
		{
			Name: "composite_types",
			PkgDir: "test_packages/composite_types",
			PkgImportPath: "some/import/composite_types",
			Expected: &gopackage.Contents{
				Functions: []gopackage.Function{
					{
						Name: "SomePointerFunc",
						Import: "some/import/composite_types",
						Args: []gopackage.Type{
							gopackage.TypePointer{
								ValueType: gopackage.TypeNamed{Name: "float32"},
							},
							gopackage.TypePointer{
								ValueType: gopackage.TypeNamed{
									Name: "SomePointerStruct",
									Import: "some/import/composite_types",
								},
							},
						},
						ReturnArgs: []gopackage.Type{
							gopackage.TypePointer{
								ValueType: gopackage.TypeNamed{Name: "string"},
							},
						},
					},
					{
						Name: "SomeArrayFunc",
						Import: "some/import/composite_types",
						Args: []gopackage.Type{
							gopackage.TypeArray{
								ValueType: gopackage.TypeNamed{
									Name: "Decimal",
									Import: "github.com/shopspring/decimal",
								},
							},
							gopackage.TypeArray{
								ValueType: gopackage.TypeNamed{Name: "float32"},
							},
						},
						ReturnArgs: []gopackage.Type{
							gopackage.TypeArray{
								ValueType: gopackage.TypeNamed{
									Name: "SomeArrayStruct",
									Import: "some/import/composite_types",
								},
							},
						},
					},
				},
				StructTypes: []gopackage.StructDecl{
					{
						Name: "SomeArrayStruct",
						Import: "some/import/composite_types",
						Fields: map[string]gopackage.Type{
							"AOfInts": gopackage.TypeArray{
								ValueType: gopackage.TypeNamed{Name: "int64"},
							},
							"AOfPToStrings": gopackage.TypeArray{
								ValueType: gopackage.TypePointer{
									ValueType: gopackage.TypeNamed{Name: "string"},
								},
							},
						},
					},
					{
						Name: "SomePointerStruct",
						Import: "some/import/composite_types",
						Fields: map[string]gopackage.Type{
							"PToInt": gopackage.TypePointer{
								ValueType: gopackage.TypeNamed{Name: "int32"},
							},
						},
					},
				},
			},
		},
		{
			Name: "proto_conversion_package",
			PkgDir: "test_packages/proto_conversion",
			PkgImportPath: "some/import/proto_conversion",
			Expected: &gopackage.Contents{
				Functions: []gopackage.Function{
					{
						Name: "IntAsStringFromProto",
						Import: "some/import/proto_conversion",
						Args: []gopackage.Type{
							gopackage.TypePointer{
								ValueType: gopackage.TypeNamed{
									Name: "IntAsString",
									Import: "some/import/proto_conversion",
								},
							},
						},
						ReturnArgs: []gopackage.Type{
							gopackage.TypeNamed{Name: "int"},
							gopackage.TypeNamed{Name: "error"},
						},
					},
					{
						Name: "IntAsStringToProto",
						Import: "some/import/proto_conversion",
						Args: []gopackage.Type{
							gopackage.TypeNamed{Name: "int"},
						},
						ReturnArgs: []gopackage.Type{
							gopackage.TypePointer{
								ValueType: gopackage.TypeNamed{
									Name: "IntAsString",
									Import: "some/import/proto_conversion",
								},
							},
							gopackage.TypeNamed{Name: "error"},
						},
					},
					{
						Name: "ShopspringDecimalFromProto",
						Import: "some/import/proto_conversion",
						Args: []gopackage.Type{
							gopackage.TypePointer{
								ValueType: gopackage.TypeNamed{
									Name: "ShopspringDecimal",
									Import: "some/import/proto_conversion",
								},
							},
						},
						ReturnArgs: []gopackage.Type{
							gopackage.TypeNamed{
								Name: "Decimal",
								Import: "github.com/shopspring/decimal",
							},
							gopackage.TypeNamed{Name: "error"},
						},
					},
					{
						Name: "ShopspringDecimalToProto",
						Import: "some/import/proto_conversion",
						Args: []gopackage.Type{
							gopackage.TypeNamed{
								Name: "Decimal",
								Import: "github.com/shopspring/decimal",
							},
						},
						ReturnArgs: []gopackage.Type{
							gopackage.TypePointer{
								ValueType: gopackage.TypeNamed{
									Name: "ShopspringDecimal",
									Import: "some/import/proto_conversion",
								},
							},
							gopackage.TypeNamed{Name: "error"},
						},
					},
					{
						Name: "init",
						Import: "some/import/proto_conversion",
					},
					{
						Name: "init",
						Import: "some/import/proto_conversion",
					},
				},
				StructTypes: []gopackage.StructDecl{
					{
						Name: "IntAsString",
						Import: "some/import/proto_conversion",
						Fields: map[string]gopackage.Type{
							"Value": gopackage.TypeNamed{Name: "string"},
							"XXX_NoUnkeyedLiteral": gopackage.TypeNamed{Name: "struct{}"},
							"XXX_unrecognized": gopackage.TypeArray{
								ValueType: gopackage.TypeNamed{Name: "byte"},
							},
							"XXX_sizecache": gopackage.TypeNamed{Name: "int32"},
						},
					},
					{
						Name: "ShopspringDecimal",
						Import: "some/import/proto_conversion",
						Fields: map[string]gopackage.Type{
							"Value": gopackage.TypeNamed{Name: "string"},
							"XXX_NoUnkeyedLiteral": gopackage.TypeNamed{Name: "struct{}"},
							"XXX_unrecognized": gopackage.TypeArray{
								ValueType: gopackage.TypeNamed{Name: "byte"},
							},
							"XXX_sizecache": gopackage.TypeNamed{Name: "int32"},
						},
					},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			pc, err := gopackage.GetContents(test.PkgDir, test.PkgImportPath)
			require.NoError(t, err)

			assert.ElementsMatch(t, test.Expected.Functions, pc.Functions)
			assert.ElementsMatch(t, test.Expected.StructTypes, pc.StructTypes)
		})
	}
}
