package gopackage

import(
	"fmt"
	"go/parser"
	"go/token"
	"go/ast"
	"path"

	"reflect"
)

type Contents struct {
	Functions []Function
	StructTypes []string
}

type Function struct {
	Name string
	Args []Type
	ReturnArgs []Type
}

type Type struct {
	Name string
	Import string
}

const CURRENT_PKG = "current_pkg_import"

func GetContents(pkgDir string, pkgImportPath string) (*Contents, error) {

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

func parseAst(pkgImportPath string, p *ast.Package) (*Contents, error) {

	pc := Contents{
		Functions: make([]Function, 0),
		StructTypes: make([]string, 0),
	}

	currentFileImports := make(map[string]string)

	ast.Inspect(p, func(node ast.Node) bool {

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

				//fmt.Println("func:", n.Name)

				f := Function{
					Name: n.Name.String(),
					Args: getTypeList(currentFileImports, n.Type.Params),
					ReturnArgs: getTypeList(currentFileImports, n.Type.Results),
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
							_, importPrefix := path.Split(pkgImportPath)
							pc.StructTypes = append(
								pc.StructTypes,
								importPrefix + "." + s.Name.Name,
							)
						}
				}
				return true

			default:
				return true
		}
	})

	return &pc, nil
}

func getTypeList(imports map[string]string, fieldList *ast.FieldList) []Type {

	if fieldList == nil || fieldList.List == nil {
		return nil
	}

	typeList := make([]Type, 0, len(fieldList.List))

	for i := range fieldList.List {
		imp, fullType := getImportAndFullType(imports, fieldList.List[i].Type)
		fmt.Println("Arg:", imp, fullType)

		typeList = append(typeList, Type{
			Name: fullType,
			Import: imp,
		})
	}

	return typeList
}

func getImportAndFullType(
	imports map[string]string,
	t ast.Expr,
) (string, string) {

	fmt.Println("******", reflect.TypeOf(t))

	switch t := t.(type) {
		case *ast.ArrayType:
			if t.Len != nil {
				panic("[...]T array types not supported")
			}
			imp, fullType := getImportAndFullType(imports, t.Elt)
			return imp, "[]" + fullType

		case *ast.Ident:
			if isBuiltInType(t.Name) {
				return "", t.Name
			}

			importPath := imports[CURRENT_PKG]
			_, importPrefix := path.Split(importPath)
			return importPath, importPrefix + "." + t.Name

		case *ast.StarExpr:
			imp, fullType := getImportAndFullType(imports, t.X)
			return imp, "*" + fullType

		case *ast.SelectorExpr:
			imp, ok := t.X.(*ast.Ident)

			if !ok {
				panic("uknown selector X")
			}

			localImport := imp.Name
			if localImport == "" {
				localImport = CURRENT_PKG
			}

			importPath, ok := imports[localImport]
			if !ok {
				panic("unknown import path " + localImport)
			}
			_, importPrefix := path.Split(importPath)

			//fmt.Println("****** Type:", importPath, importPrefix + "." + t.Sel.Name)

			return importPath, importPrefix + "." + t.Sel.Name

		default:
			panic("unknown field type")
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
		"error": {},
		"int": {},
	}

	_, ok := builtInTypes[t]
	return ok
}
