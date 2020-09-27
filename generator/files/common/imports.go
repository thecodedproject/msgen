package common

import(
	"github.com/thecodedproject/msgen/parser"
	"github.com/thecodedproject/msgen/parser/proto_helpers"
	"path"
	"sort"
)

func ServiceNameFromRootImportPath(rootImportPath string) string {

	_, serviceName := path.Split(rootImportPath)
	return serviceName
}

func SortedImportsWithNestedTypesImport(
	serviceRootImportPath string,
	i parser.ProtoInterface,
	extraImports... string,
) []string {

	imports := make([]string, 0, len(extraImports) + 1)

	nestedMessages := proto_helpers.NestedMessages(i)
	if len(nestedMessages) > 0 {
		imports = append(imports, serviceRootImportPath)
	}

	imports = append(imports, extraImports...)

	sort.Strings(imports)

	return imports
}
