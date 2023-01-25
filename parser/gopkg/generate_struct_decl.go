package gopkg

import (
	"errors"
	"io"
)

func GenerateStructDecl(
	w io.Writer,
	decl DeclStruct,
	importAliases map[string]string,
) error {

	if decl.Name == "" {
		return errors.New("struct name cannot be empty")
	}

	w.Write([]byte("type MyObject struct{}\n"))

	return nil
}
