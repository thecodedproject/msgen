package files

import(
	"github.com/iancoleman/strcase"
	"github.com/thecodedproject/msgen/parser"
	"strings"
	"text/template"
)

func BaseTemplate() *template.Template {

	return template.New("").Funcs(map[string]interface{}{
		"FuncRetVals": funcReturnSignature(false, false),
		"NamedFuncRetVals": funcReturnSignature(true, false),
		"NamedFuncRetValsWithError": funcReturnSignature(true, true),

		"FuncDefaultReturn_Named_WithError": funcDefaultReturnStatement(true, true),

		"ToLower": strings.ToLower,
		"ToCamel": strcase.ToCamel,
		"ToLowerCamel": strcase.ToLowerCamel,
	})
}

func funcReturnSignature(namedReturnVariables, addError bool) func([]parser.Field) string {

	return func(returnArgs []parser.Field) string {

		shouldBraceReturnSignature := len(returnArgs) > 1 || namedReturnVariables || addError

		var retSignature string
		if shouldBraceReturnSignature {
			retSignature = "("
		}

		for i := range returnArgs {
			if namedReturnVariables {
				retSignature += strcase.ToLowerCamel("res_" + returnArgs[i].Name) + " " + returnArgs[i].Type
			} else {
				retSignature += returnArgs[i].Type
			}

			if i != len(returnArgs)-1 {
				retSignature += ", "
			}
		}

		if addError {
			if namedReturnVariables {
				retSignature += ", err error)"
			} else {
				retSignature += ", error)"
			}
		} else if shouldBraceReturnSignature {
			retSignature += ")"
		}

		return retSignature
	}
}

func funcDefaultReturnStatement(namedReturnVariables, addError bool) func([]parser.Field) string {

	if !namedReturnVariables {
		panic("Unamed funcDefaultReturnStatement not implemented")
	}

	return func(returnArgs []parser.Field) string {

		returnStatement := "return "

		for i := range returnArgs {
			returnStatement += strcase.ToLowerCamel("res_" + returnArgs[i].Name)
			if i != len(returnArgs)-1 {
				returnStatement += ", "
			}
		}

		if addError {
			returnStatement += ", err"
		}

		return returnStatement
	}
}
