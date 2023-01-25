package gopkg

import (
)

type DeclFunc struct {
	Name string
	Import string
	Args []DeclVar
	ReturnArgs []Type
	BodyTmpl string
}

type DeclStruct struct {
	Name string
	Import string
	// TODO: Maybe also add struct field descriptors
	Fields []DeclVar//map[string]Type
}

type DeclType struct {
	Name string
	Import string
	Type Type
}

type DeclVar struct {
	Type
	Name string
	Import string
}

func (f DeclFunc) FullDecl(importAliases map[string]string) string {

	decl := "func " + f.Name + "("

	if len(f.Args) == 1 {
		decl += f.Args[0].Name + " " + f.Args[0].FullType(importAliases)
	} else if len(f.Args) > 1 {
		decl += "\n"
		for _, arg := range f.Args {
			decl += "\t" +  arg.Name + " " + arg.FullType(importAliases) + ",\n"
		}
	}

	decl += ")"

	if len(f.ReturnArgs) == 1 {
			decl += " " + f.ReturnArgs[0].FullType(importAliases)
	} else if len(f.ReturnArgs) > 1 {
		decl += " ("
		for i, ret := range f.ReturnArgs {
			decl += ret.FullType(importAliases)

			if i < len(f.ReturnArgs) - 1 {
				decl += ", "
			}
		}
		decl += ")"
	}

	return decl
}
