package generator

import(
	"github.com/pkg/errors"
	"github.com/thecodedproject/msgen/generator/files/common"
	"io"
	"io/ioutil"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"sort"
)

func Generate(
	inputFile string,
	inputStructName string,
	outputInterfaceName string,
	writer io.Writer,
) error {

	stateInfo, err := parseInputFile(inputFile, inputStructName)
	if err != nil {
		return err
	}

	stateInfo = removeImportsNotUsedForFields(stateInfo)

	stateInfo.Imports = append(
		stateInfo.Imports,
		Import{Path: "\"testing\""},
	)

	sortImports(stateInfo.Imports)

	t, err := common.BaseTemplate().Parse(backendsTmpl)
	if err != nil {
		return err
	}

	return t.Execute(writer, stateInfo)
}

type Field struct {
	Name string
	Type string
}

type Import struct {
	Path string
	Alias string
}

type StateInfo struct {
	PackageName string
	Fields []Field
	Imports []Import
}

func parseInputFile(
	inputFile string,
	inputStructName string,
) (StateInfo, error) {

	src, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return StateInfo{}, errors.Wrap(err, "Error reading intput file")
	}

	astFile, err := parser.ParseFile(token.NewFileSet(), "", src, 0)
	if err != nil {
		return StateInfo{}, errors.Wrap(err, "Error parsing input file")
	}

	var info StateInfo
	var inspectErr error
	ast.Inspect(astFile, func(node ast.Node) bool {

		switch n := node.(type) {

		case *ast.File:
			info.PackageName = n.Name.Name

		case *ast.ImportSpec:
			i := Import{
				Path:  n.Path.Value,
			}
			if n.Name != nil {
				i.Alias = n.Name.Name
			}
			info.Imports = append(info.Imports, i)

		case *ast.GenDecl:
			switch s := n.Specs[0].(type) {
			case *ast.TypeSpec:
				if s.Name.Name == inputStructName {

					switch inputType := s.Type.(type) {
					case *ast.StructType:

						if inputType.Fields != nil {
							for _, f := range inputType.Fields.List {


								if len(f.Names) != 1 || f.Names[0] == nil {
									inspectErr = errors.New("input type has field with no name")
									return false
								}

								info.Fields = append(info.Fields, Field{
									Name: f.Names[0].Name,
									Type: string(src[f.Type.Pos()-1:f.Type.End()-1]),
								})

							}

						}


					default:
						inspectErr = errors.New("input type is not struct type")
						return false
					}

					return false
				}

			}
		}

		return true
	})

	if inspectErr != nil {
		return StateInfo{}, inspectErr
	}

	return info, nil
}

func removeImportsNotUsedForFields(inputStructInfo StateInfo) StateInfo {

	var requiredImports []Import

	for _, f := range inputStructInfo.Fields {

		fieldsTypeParts := strings.Split(f.Type, ".")

		if len(fieldsTypeParts) == 1 {
			continue
		}

		fieldPackageName := fieldsTypeParts[0]

		for _, i := range inputStructInfo.Imports {

			var packageName string
			if i.Alias != "" {
				packageName = i.Alias
			} else {
				pathPackages := strings.Split(i.Path, "/")
				packageName = strings.Replace(
					pathPackages[len(pathPackages)-1],
					"\"",
					"",
					-1,
				)
			}

			if fieldPackageName == packageName {
				requiredImports = append(requiredImports, i)
			}
		}
	}

	inputStructInfo.Imports = requiredImports

	return inputStructInfo
}

func sortImports(imports []Import) {

	sort.Slice(imports, func(i, j int) bool {
		if imports[i].Alias == imports[j].Alias {
			return imports[i].Path < imports[j].Path
		}
		return imports[i].Alias < imports[j].Alias
	})
}

var backendsTmpl = `package state

imports(
{{- range .Imports}}
	{{if .Alias}}{{.Alias}} {{end}}{{.Path}}
{{- end}}
)

type State interface {
{{- range .Fields}}
	Get{{ToCamel .Name}}() {{.Type}}
{{- end}}
}

{{- range .Fields}}

func (s *stateImpl) Get{{ToCamel .Name}}() {{.Type}} {

	return s.{{ToLowerCamel .Name}}
}
{{- end}}

type stateOption func(*stateImpl)

func NewStateForTesting(
	_ testing.TB,
	opts ...stateOption,
) *State {

	var s State
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}

{{- range .Fields}}

func With{{ToCamel .Name}}({{ToLowerCamel .Name}} {{.Type}}) stateOption {

	return func(s *stateImpl) {
		s.{{ToLowerCamel .Name}} = {{ToLowerCamel .Name}}
	}
}
{{- end}}
`
