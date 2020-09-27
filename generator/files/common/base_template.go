package common

import(
	"github.com/iancoleman/strcase"
	"github.com/thecodedproject/msgen/parser"
	"strings"
	"text/template"
)

func BaseTemplate(nestedTypeImport string) *template.Template {

	return template.New("").Funcs(map[string]interface{}{
		"FuncArgs": funcArgumentSignature(nestedTypeImport, false, false, false),
		"FuncArgsWithCtx": funcArgumentSignature(nestedTypeImport, false, true, false),
		"FuncArgsWithCtxAndState": funcArgumentSignature(nestedTypeImport, false, true, true),
		"SplitFuncArgs": funcArgumentSignature(nestedTypeImport, true, false, false),
		"SplitFuncArgsWithCtx": funcArgumentSignature(nestedTypeImport, true, true, false),
		"SplitFuncArgsWithCtxAndState": funcArgumentSignature(nestedTypeImport, true, true, true),

		"FuncRetVals": funcReturnSignature(nestedTypeImport, false, false),
		"FuncRetValsWithError": funcReturnSignature(nestedTypeImport, false, true),
		"NamedFuncRetVals": funcReturnSignature(nestedTypeImport, true, false),
		"NamedFuncRetValsWithError": funcReturnSignature(nestedTypeImport, true, true),

		"FuncDefaultReturn_WithError": funcDefaultReturnStatement(nestedTypeImport, false, true),
		"FuncDefaultReturn_Named_WithError": funcDefaultReturnStatement(nestedTypeImport, true, true),

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

func funcArgumentSignature(nestedTypeImport string, splitLines, withContext, withState bool) func([]parser.Field) string {

	return func(args []parser.Field) string {

		if withState {
			args = append([]parser.Field{{Name: "s", Type: "state.State"}}, args...)
		}
		if withContext {
			args = append([]parser.Field{{Name: "ctx", Type: "context.Context"}}, args...)
		}

		argSig := "("
		for i, f := range args {
			argType := buildArgType(nestedTypeImport, f)

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

func funcReturnSignature(nestedTypeImport string, namedReturnVariables, addError bool) func([]parser.Field) string {

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

			returnType := buildArgType(nestedTypeImport, returnArgs[i])

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

func buildArgType(nestedTypeImport string, f parser.Field) string {

	if !f.IsNestedMessage {
		return f.Type
	}

	if nestedTypeImport == "" {
		return "*" + f.Type
	}

	return "*" + nestedTypeImport + "." + f.Type
}

func funcDefaultReturnStatement(nestedTypeImport string, namedReturnVariables, addError bool) func([]parser.Field) string {

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

		for i, arg := range returnArgs {
			if namedReturnVariables {
				returnStatement += strcase.ToLowerCamel("res_" + arg.Name)
			} else if arg.IsNestedMessage {
				if nestedTypeImport == "" {
					returnStatement += "&" + arg.Type + "{}"
				} else {
					returnStatement += "&" + nestedTypeImport + "." + arg.Type + "{}"
				}
			} else {
				defaultRetVal, ok := defaultReturnValues[arg.Type]
				if !ok {
					panic("No default return value for type " + arg.Type)
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
