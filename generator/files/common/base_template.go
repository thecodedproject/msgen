package common

import(
	"github.com/iancoleman/strcase"
	"github.com/thecodedproject/msgen/parser"
	"strings"
	"text/template"
)

func BaseTemplate() *template.Template {

	return template.New("").Funcs(map[string]interface{}{
		"FuncRetVals": funcReturnSignature(false, false),
		"FuncRetValsWithError": funcReturnSignature(false, true),
		"NamedFuncRetVals": funcReturnSignature(true, false),
		"NamedFuncRetValsWithError": funcReturnSignature(true, true),

		"FuncDefaultReturn_WithError": funcDefaultReturnStatement(false, true),
		"FuncDefaultReturn_Named_WithError": funcDefaultReturnStatement(true, true),

		"ToLower": strings.ToLower,
		"ToCamel": strcase.ToCamel,
		"ToLowerCamel": strcase.ToLowerCamel,
	})
}

// TODO(jcooper): Make this more extensible by storing it in a JSON file
var defaultReturnValues = map[string]string{
	"error": "nil",
	"float32": "0",
	"float64": "0",
	"int": "0",
	"int32": "0",
	"int64": "0",
	"string": "\"\"",
}

func funcReturnSignature(namedReturnVariables, addError bool) func([]parser.Field) string {

	return func(returnArgs []parser.Field) string {

		if len(returnArgs) == 0 {
			if !addError {
				return ""
			}

			if namedReturnVariables {
				return "(err error)"
			}
			return "error"
		}

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

	return func(returnArgs []parser.Field) string {

		if len(returnArgs) == 0 {
			if !addError {
				return "return"
			}

			if namedReturnVariables {
				return "return err"
			}

			return "return " + defaultReturnValues["error"]
		}

		returnStatement := "return "

		for i := range returnArgs {
			if namedReturnVariables {
				returnStatement += strcase.ToLowerCamel("res_" + returnArgs[i].Name)
			} else {
				defaultRetVal, ok := defaultReturnValues[returnArgs[i].Type]
				if !ok {
					panic("No default return value for type" + returnArgs[i].Type)
				}
				returnStatement += defaultRetVal
			}

			if i != len(returnArgs)-1 {
				returnStatement += ", "
			}
		}

		if addError {
			returnStatement += ", "
			if namedReturnVariables {
				returnStatement += "err"
			} else {
				returnStatement += defaultReturnValues["error"]
			}

		}

		return returnStatement
	}
}
