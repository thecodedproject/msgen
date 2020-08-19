package common

import(
	"path"
	"os"
)

func CreatePathAndOpen(
	filepath string,
) (*os.File, error) {

	dir, _ := path.Split(filepath)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return os.Create(filepath)
}
