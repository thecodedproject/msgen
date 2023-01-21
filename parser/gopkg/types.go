package gopkg

import (
	"errors"
)

type Contents struct {
	Functions []DeclFunc
	StructTypes []DeclStruct
}

type Type interface {
	DefaultInit(importAliases map[string]string) (string, error)
	FullType(importAliases map[string]string) string
}

type TypeArray struct {
	ValueType Type
}

func (t TypeArray) DefaultInit(importAliases map[string]string) (string, error) {
	return "nil", nil
}

func (t TypeArray) FullType(importAliases map[string]string) string {
	return "[]" + t.ValueType.FullType(importAliases)
}

type TypeByte struct {}

func (t TypeByte) DefaultInit(importAliases map[string]string) (string, error) {
	return "0", nil
}

func (t TypeByte) FullType(importAliases map[string]string) string {
	return "byte"
}

type TypeError struct {}

func (t TypeError) DefaultInit(importAliases map[string]string) (string, error) {
	return "nil", nil
}

func (t TypeError) FullType(importAliases map[string]string) string {
	return "error"
}

type TypeFloat32 struct {}

func (t TypeFloat32) DefaultInit(importAliases map[string]string) (string, error) {
	return "0", nil
}

func (t TypeFloat32) FullType(importAliases map[string]string) string {
	return "float32"
}

type TypeFloat64 struct {}

func (t TypeFloat64) DefaultInit(importAliases map[string]string) (string, error) {
	return "0", nil
}

func (t TypeFloat64) FullType(importAliases map[string]string) string {
	return "float64"
}

type TypeInt struct {}

func (t TypeInt) DefaultInit(importAliases map[string]string) (string, error) {
	return "0", nil
}

func (t TypeInt) FullType(importAliases map[string]string) string {
	return "int"
}

type TypeInterface struct {
	Name string
	Import string
	// TODO maybe add funcs?
}

func (t TypeInterface) DefaultInit(importAliases map[string]string) (string, error) {
	return "nil", nil
}

func (t TypeInterface) FullType(importAliases map[string]string) string {
	if alias, hasAlias := importAliases[t.Import]; hasAlias {
		return alias + "." + t.Name
	}

	return t.Name
}

type TypeInt32 struct {}

func (t TypeInt32) DefaultInit(importAliases map[string]string) (string, error) {
	return "0", nil
}

func (t TypeInt32) FullType(importAliases map[string]string) string {
	return "int32"
}

type TypeInt64 struct {}

func (t TypeInt64) DefaultInit(importAliases map[string]string) (string, error) {
	return "0", nil
}

func (t TypeInt64) FullType(importAliases map[string]string) string {
	return "int64"
}

type TypeString struct {}

func (t TypeString) DefaultInit(importAliases map[string]string) (string, error) {
	return "\"\"", nil
}

func (t TypeString) FullType(importAliases map[string]string) string {
	return "string"
}

type TypeStruct struct {
	Name string
	Import string
	// TODO add fields
}

func (t TypeStruct) DefaultInit(importAliases map[string]string) (string, error) {
	return t.FullType(importAliases) + "{}", nil
}

func (t TypeStruct) FullType(importAliases map[string]string) string {
	if alias, hasAlias := importAliases[t.Import]; hasAlias {
		return alias + "." + t.Name
	}

	return t.Name
}

// TypeUnknownNamed represents a named type who exact type is unknown
// This occurs typically when parsing a package which imports a type
// from another package - it's not known at parsing time whether the other
// type is a struct, interface, typedef or a variable.
type TypeUnknownNamed struct {
	Name string
	Import string
}

func (t TypeUnknownNamed) DefaultInit(importAliases map[string]string) (string, error) {
	// TODO return more informative error
	return "...", errors.New("not implemented")
}

func (t TypeUnknownNamed) FullType(importAliases map[string]string) string {
	if alias, hasAlias := importAliases[t.Import]; hasAlias {
		return alias + "." + t.Name
	}

	return t.Name
}

type TypeMap struct {
	KeyType Type
	ValueType Type
}

func (t TypeMap) DefaultInit(importAliases map[string]string) (string, error) {
	return "nil", nil
}

func (t TypeMap) FullType(importAliases map[string]string) string {
	return "map[" + t.KeyType.FullType(importAliases) + "]" + t.ValueType.FullType(importAliases)
}

type TypePointer struct {
	ValueType Type
}

func (t TypePointer) DefaultInit(importAliases map[string]string) (string, error) {
	return "nil", nil
}

func (t TypePointer) FullType(importAliases map[string]string) string {
	return "*" + t.ValueType.FullType(importAliases)
}
