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
			Name: "some",
			PkgDir: "test_packages/proto_conversion",
			PkgImportPath: "some/import/proto_conversion",
			Expected: &gopackage.Contents{
				Functions: []gopackage.Function{
					{
						Name: "IntAsStringFromProto",
						Args: []gopackage.Type{
							{
								Name: "*proto_conversion.IntAsString",
								Import: "some/import/proto_conversion",
							},
						},
						ReturnArgs: []gopackage.Type{
							{
								Name: "int",
							},
							{
								Name: "error",
							},
						},
					},
					{
						Name: "IntAsStringToProto",
						Args: []gopackage.Type{
							{
								Name: "int",
							},
						},
						ReturnArgs: []gopackage.Type{
							{
								Name: "*proto_conversion.IntAsString",
								Import: "some/import/proto_conversion",
							},
							{
								Name: "error",
							},
						},
					},
					{
						Name: "ShopspringDecimalFromProto",
						Args: []gopackage.Type{
							{
								Name: "*proto_conversion.ShopspringDecimal",
								Import: "some/import/proto_conversion",
							},
						},
						ReturnArgs: []gopackage.Type{
							{
								Name: "decimal.Decimal",
								Import: "github.com/shopspring/decimal",
							},
							{
								Name: "error",
							},
						},
					},
					{
						Name: "ShopspringDecimalToProto",
						Args: []gopackage.Type{
							{
								Name: "decimal.Decimal",
								Import: "github.com/shopspring/decimal",
							},
						},
						ReturnArgs: []gopackage.Type{
							{
								Name: "*proto_conversion.ShopspringDecimal",
								Import: "some/import/proto_conversion",
							},
							{
								Name: "error",
							},
						},
					},
					{
						Name: "init",
					},
					{
						Name: "init",
					},
				},
				StructTypes: []string{
					"proto_conversion.IntAsString",
					"proto_conversion.ShopspringDecimal",
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
