package parse

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"go/ast"
	"path"

	"github.com/thecodedproject/msgen/parser/gopkg"

	//"reflect"
)

const CURRENT_PKG = "current_pkg_import"

func GetContents(pkgDir string, pkgImportPath string) (*gopkg.Contents, error) {

	pkgs, err := parser.ParseDir(
		token.NewFileSet(),
		pkgDir,
		nil,
		0,
	)
	if err != nil {
		return nil, err
	}

	if len(pkgs) != 1 {
		for k := range pkgs {
			fmt.Println(k)
		}

		return nil, fmt.Errorf("more than one package found in dir %s", pkgDir)
	}

	for k := range pkgs {
		return parseAst(pkgImportPath, pkgs[k])
	}

	return nil, nil
}

func parseAst(pkgImportPath string, p *ast.Package) (*gopkg.Contents, error) {

	pc := gopkg.Contents{
		Functions: make([]gopkg.DeclFunc, 0),
		StructTypes: make([]gopkg.DeclStruct, 0),
	}

	currentFileImports := make(map[string]string)

	var inspectingErr error
	ast.Inspect(p, func(node ast.Node) bool {

		// If we have encountered an error stop parsing the AST asap (by stopping
		// any more recursion into the ast)
		if inspectingErr != nil {
			return false
		}

		switch n := node.(type) {

			case *ast.ImportSpec:
				addImport(currentFileImports, n)
				return true

			case *ast.File:

				//fmt.Println("File: ", n.Name)

				//fmt.Printf("Previous imports: %+v\n", currentFileImports)
				currentFileImports = make(map[string]string)
				currentFileImports[CURRENT_PKG] = pkgImportPath
				return true

			case *ast.FuncDecl:

				// Skip functions with receivers
				if n.Recv != nil {
					return true
				}

				args, err := getDeclVarsFromFieldList(currentFileImports, n.Type.Params)
				if err != nil {
					inspectingErr = err
					return false
				}

				retArgs, err := getArgTypeList(currentFileImports, n.Type.Results)
				if err != nil {
					inspectingErr = err
					return false
				}

				f := gopkg.DeclFunc{
					Name: n.Name.String(),
					Import: pkgImportPath,
					Args: args,
					ReturnArgs: retArgs,
				}

				pc.Functions = append(pc.Functions, f)
				return true

			case *ast.GenDecl:
				if len(n.Specs) == 0 {
					return true
				}

				switch s := n.Specs[0].(type) {
				case *ast.TypeSpec:
						switch s.Type.(type) {
						case *ast.StructType:

							structSpec := s.Type.(*ast.StructType)

							structFields, err := getFieldTypeList(
								currentFileImports,
								structSpec.Fields,
							)
							if err != nil {
								inspectingErr = err
								return false
							}

							pc.StructTypes = append(
								pc.StructTypes,
								gopkg.DeclStruct{
									Name: s.Name.Name,
									Import: pkgImportPath,
									Fields: structFields,
								},
							)
						}
				}
				return true

			default:
				return true
		}
	})

	if inspectingErr != nil {
		return nil, inspectingErr
	}

	return &pc, nil
}

// getArgTypeList gets an order list of arguments from an `ast.FieldList`
//
// Used to get either the types of the parameters arguments for a function,
// or the return arguments for a function whilst parsing the ast.
func getArgTypeList(
	imports map[string]string,
	fieldList *ast.FieldList,
) ([]gopkg.Type, error) {

	if fieldList == nil || fieldList.List == nil {
		return nil, nil
	}

	typeList := make([]gopkg.Type, 0, len(fieldList.List))

	for i := range fieldList.List {
		fieldType, err := getFullType(imports, fieldList.List[i].Type)
		if err != nil {
			return nil, err
		}
		typeList = append(typeList, fieldType)
	}

	return typeList, nil
}

// getDeclVarsFromFieldList returns an ordered list of declared variables
//
// The ast field list might be, for example, the list of arguments passed into
// a function.
// It returns the underlying type (as `gopkg.Type`) as well as the name of the
// declared variable.
// Note that `gopkg.DeclVar.Import` will always be blank as the field list will
// only contain vars declared in a local scope (i.e. not at the package level)
func getDeclVarsFromFieldList(
	imports map[string]string,
	fieldList *ast.FieldList,
) ([]gopkg.DeclVar, error) {

	if fieldList == nil || fieldList.List == nil {
		return nil, nil
	}

	typeList := make([]gopkg.DeclVar, 0, len(fieldList.List))

	for i := range fieldList.List {
		fieldType, err := getFullType(imports, fieldList.List[i].Type)
		if err != nil {
			return nil, err
		}
		typeList = append(typeList, gopkg.DeclVar{
			Name: fieldList.List[i].Names[0].String(),
			Type: fieldType,
		})
	}

	return typeList, nil
}


// getFieldTypeList returns a map of field names and types from an `ast.FieldList`
//
// Used to get the list of fields and there types when pasing the ast for a struct
func getFieldTypeList(
	imports map[string]string,
	fieldList *ast.FieldList,
) (map[string]gopkg.Type, error) {

	if fieldList == nil || fieldList.List == nil {
		return nil, nil
	}

	fieldTypeList := make(map[string]gopkg.Type)

	for i := range fieldList.List {
		fieldType, err := getFullType(imports, fieldList.List[i].Type)
		if err != nil {
			return nil, err
		}
		fieldTypeList[fieldList.List[i].Names[0].String()] = fieldType
	}

	return fieldTypeList, nil
}

func getFullType(
	imports map[string]string,
	t ast.Expr,
) (gopkg.Type, error) {

	//fmt.Println("******", reflect.TypeOf(t))

	switch t := t.(type) {
		case *ast.ArrayType:
			if t.Len != nil {
				return nil, errors.New("[...]T array types not supported")
			}
			fullType, err := getFullType(imports, t.Elt)
			if err != nil {
				return nil, err
			}
			return gopkg.TypeArray{
				ValueType: fullType,
			}, nil

		case *ast.Ident:
			if isBuiltInType(t.Name) {

				return typeFromString(t.Name), nil

				//return gopkg.TypeUnknownNamed{
				//	Name: t.Name,
				//}, nil
			}

			importPath := imports[CURRENT_PKG]
			return gopkg.TypeUnknownNamed{
				Name: t.Name,
				Import: importPath,
			}, nil

		case *ast.StarExpr:
			fullType, err := getFullType(imports, t.X)
			if err != nil {
				return nil, err
			}
			return gopkg.TypePointer{
				ValueType: fullType,
			}, nil

		// i.e. an expression selecting something from another package
		//	`some_pkg.SomeType`
		case *ast.SelectorExpr:
			imp, ok := t.X.(*ast.Ident)

			if !ok {
				return nil, errors.New("uknown selector X")
			}

			importPath, ok := imports[imp.Name]
			if !ok {
				return nil, errors.New("unknown import path " + imp.Name)
			}

			//fmt.Println("****** Type:", importPath, importPrefix + "." + t.Sel.Name)

			return gopkg.TypeUnknownNamed{
				Name: t.Sel.Name,
				Import: importPath,
			}, nil

		case *ast.StructType:
			return gopkg.TypeUnknownNamed{
				Name: "struct{}",
			}, nil

		default:
			return nil, errors.New("unknown field type")
	}
}

func removeQuotes(s string) string {

	if s[0] == '"' {
		s = s[1:]
	}
	if s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

func addImport(imports map[string]string, n *ast.ImportSpec) {

	importPath := removeQuotes(n.Path.Value)
	var localName string
	if n.Name != nil {
		localName = n.Name.String()
	}
	if localName == "." {
		panic("'.' imports are not supported")
	}
	if localName == "" {
		_, localName = path.Split(importPath)
	}

	imports[localName] = importPath
}

func isBuiltInType(t string) bool {

	builtInTypes := map[string]struct{}{
		"byte": {},
		"error": {},
		"float32": {},
		"float64": {},
		"int": {},
		"int32": {},
		"int64": {},
		"string": {},
	}

	_, ok := builtInTypes[t]
	return ok
}

func typeFromString(t string) gopkg.Type {

	switch t {
	case "byte":
		return gopkg.TypeByte{}
	case "error":
		return gopkg.TypeError{}
	case "float32":
		return gopkg.TypeFloat32{}
	case "float64":
		return gopkg.TypeFloat64{}
	case "int":
		return gopkg.TypeInt{}
	case "int32":
		return gopkg.TypeInt32{}
	case "int64":
		return gopkg.TypeInt64{}
	case "string":
		return gopkg.TypeString{}
	}
	return nil
}
