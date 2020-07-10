package proto_helpers

import(
	"fmt"
	"github.com/thecodedproject/msgen/parser"
)

// This file contins helper functions for accessing/manipulating a ProtoInterface

func Method(
	i parser.ProtoInterface,
	methodName string,
) (parser.Method, error) {

	for _, m := range i.Methods {
		if m.Name == methodName {
			return m, nil
		}
	}
	return parser.Method{}, fmt.Errorf("No such method '%s'", methodName)
}

func Message(
	i parser.ProtoInterface,
	messageName string,
) (parser.Message, error) {

	for _, m := range i.Messages {
		if m.Name == messageName {
			return m, nil
		}
	}
	return parser.Message{}, fmt.Errorf("No such message '%s'", messageName)
}

func MethodRequestMessage(
	i parser.ProtoInterface,
	methodName string,
) (parser.Message, error) {

	method, err := Method(i, methodName)
	if err != nil {
		return parser.Message{}, err
	}

	message, err := Message(i, method.RequestMessage)
	if err != nil {
		return parser.Message{}, fmt.Errorf("No such request message for method '%+v'", method)
	}

	return message, nil
}

func MethodResponseMessage(
	i parser.ProtoInterface,
	methodName string,
) (parser.Message, error) {

	method, err := Method(i, methodName)
	if err != nil {
		return parser.Message{}, err
	}

	message, err := Message(i, method.ResponseMessage)
	if err != nil {
		return parser.Message{}, fmt.Errorf("No such response message for method '%+v'", method)
	}

	return message, nil
}

func MethodRequestFields(
	i parser.ProtoInterface,
	methodName string,
) ([]parser.Field, error) {

	message, err := MethodRequestMessage(i, methodName)
	if err != nil {
		return nil, err
	}

	return message.Fields, nil
}

func MethodResponseFields(
	i parser.ProtoInterface,
	methodName string,
) ([]parser.Field, error) {

	message, err := MethodResponseMessage(i, methodName)
	if err != nil {
		return nil, err
	}

	return message.Fields, nil
}

func MethodResponseTypes(
	i parser.ProtoInterface,
	methodName string,
) ([]string, error) {

	message, err := MethodResponseMessage(i, methodName)
	if err != nil {
		return nil, err
	}

	return FieldTypes(message.Fields), nil
}

func FieldNames(fields []parser.Field) []string {

	types := make([]string, len(fields))
	for i := range fields {
		types[i] = fields[i].Name
	}
	return types
}

func FieldTypes(fields []parser.Field) []string {

	types := make([]string, len(fields))
	for i := range fields {
		types[i] = fields[i].Type
	}
	return types
}
