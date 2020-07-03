package files

import(
	"github.com/iancoleman/strcase"
	"strings"
	"text/template"
)

func BaseTemplate() *template.Template {

	return template.New("").Funcs(map[string]interface{}{
		"FuncRetVals": genFunctionReturn,
		"ToLower": strings.ToLower,

		"ToCamel": strcase.ToCamel,
		"ToLowerCamel": strcase.ToLowerCamel,
	})
}

func genFunctionReturn(returnTypes []string) string {

	if len(returnTypes) == 1 {
		return returnTypes[0]
	}

	returnStatement := "("

	for i := range returnTypes {
		returnStatement += returnTypes[i]

		if i != len(returnTypes)-1 {
			returnStatement += ", "
		}
	}

	returnStatement += ")"

	return returnStatement
}
