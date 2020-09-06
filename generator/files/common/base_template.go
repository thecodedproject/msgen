package common

import(
	"github.com/iancoleman/strcase"
	"github.com/thecodedproject/msgen/parser"
	"strings"
	"text/template"
)

func BaseTemplate() *template.Template {

	return template.New("").Funcs(map[string]interface{}{
		"FuncArgs": funcArgumentSignature(false, false, false),
		"FuncArgsWithCtx": funcArgumentSignature(false, true, false),
		"FuncArgsWithCtxAndState": funcArgumentSignature(false, true, true),
		"SplitFuncArgs": funcArgumentSignature(true, false, false),
		"SplitFuncArgsWithCtx": funcArgumentSignature(true, true, false),
		"SplitFuncArgsWithCtxAndState": funcArgumentSignature(true, true, true),

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

func funcArgumentSignature(splitLines, withContext, withState bool) func([]parser.Field) string {

	return func(args []parser.Field) string {

		if withState {
			args = append([]parser.Field{{Name: "s", Type: "state.State"}}, args...)
		}
		if withContext {
			args = append([]parser.Field{{Name: "ctx", Type: "context.Context"}}, args...)
		}

		argSig := "("
		for i, f := range args {
			argType := f.Type
			if f.IsNestedMessage {
				argType = "*" + argType
			}

			if splitLines {
				argSig += "\n\t"
			}

			argSig += strcase.ToLowerCamel(f.Name) + " " + argType

			if splitLines {
				argSig += ","
			} else if i != len(args)-1 {
				argSig += ", "
			}
		}
		if splitLines {
				argSig += "\n"
		}
		argSig += ")"

		return argSig
	}
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

			returnType := returnArgs[i].Type
			if returnArgs[i].IsNestedMessage {
				returnType = "*" + returnType
			}

			if namedReturnVariables {
				retSignature += strcase.ToLowerCamel("res_" + returnArgs[i].Name) + " " + returnType
			} else {
				retSignature += returnType
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
