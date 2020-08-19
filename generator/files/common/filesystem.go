package common

import(
	"path"
	"os"
)

func CreatePathAndOpen(
	filepath string,
) (*os.File, error) {

	dir, _ := path.Split(filepath)

	if dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}
	}

	return os.Create(filepath)
}

func ServiceNameFromRootImportPath(rootImportPath string) string {

	_, serviceName := path.Split(rootImportPath)
	return serviceName
}
