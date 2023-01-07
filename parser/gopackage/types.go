package gopackage

import (
)

type Contents struct {
	Functions []Function
	StructTypes []StructDecl
}

type Function struct {
	Name string
	Import string
	Args []Type
	ReturnArgs []Type
}

type StructDecl struct {
	Name string
	Import string
	// TODO: Maybe also add struct field descriptors
	Fields map[string]Type
}

type Type interface {
	FullName(importAliases map[string]string) string
}

type TypeArray struct {
	ValueType Type
}

func (t TypeArray) FullName(importAliases map[string]string) string {
	return "[]" + t.ValueType.FullName(importAliases)
}

type TypeNamed struct {
	Name string
	Import string
}

func (t TypeNamed) FullName(importAliases map[string]string) string {
	if alias, hasAlias := importAliases[t.Import]; hasAlias {
		return alias + "." + t.Name
	}

	return t.Name
}

type TypeMap struct {
	KeyType Type
	ValueType Type
}

func (t TypeMap) FullName(importAliases map[string]string) string {
	return "map[" + t.KeyType.FullName(importAliases) + "]" + t.ValueType.FullName(importAliases)
}

type TypePointer struct {
	ValueType Type
}

func (t TypePointer) FullName(importAliases map[string]string) string {
	return "*" + t.ValueType.FullName(importAliases)
}
