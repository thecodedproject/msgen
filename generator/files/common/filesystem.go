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

func FileExists(filepath string) (bool, error) {

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
