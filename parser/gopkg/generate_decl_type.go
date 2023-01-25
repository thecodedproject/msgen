package gopkg

import (
	"errors"
	"io"
)

func GenerateDeclType(
	w io.Writer,
	decl DeclType,
	importAliases map[string]string,
) error {

	if decl.Name == "" {
		return errors.New("type decl name cannot be empty")
	}
	if decl.Type == nil {
		return errors.New("type decl type cannot be nil")
	}

	w.Write([]byte("type MyStrct struct {}\n"))

	return nil
}
