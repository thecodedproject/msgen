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

func MethodRequestFieldsWithImportOnNestedFields(
	i parser.ProtoInterface,
	methodName string,
	nestedImportedName string,
) ([]parser.Field, error) {

	fields, err := MethodRequestFields(
		i,
		methodName,
	)
	if err != nil {
		return nil, err
	}

	return AddImportToNestedFieldNames(fields, nestedImportedName), nil
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

func MethodResponseFieldsWithImportOnNestedFields(
	i parser.ProtoInterface,
	methodName string,
	nestedImportedName string,
) ([]parser.Field, error) {

	fields, err := MethodResponseFields(
		i,
		methodName,
	)
	if err != nil {
		return nil, err
	}

	return AddImportToNestedFieldNames(fields, nestedImportedName), nil
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

// NestedMessages returns all of the messages that are
// nested within any other message in the proto interface
func NestedMessages(
	i parser.ProtoInterface,
) []parser.Message {

	nestedMessages := make([]parser.Message, 0)
	for _, m := range i.Messages {
		for _, f := range m.Fields {
			nestedMessage, err := Message(i, f.Type)
			if err != nil {
				continue // Not a nested type
			}
			nestedMessages = append(nestedMessages, nestedMessage)
		}
	}
	return nestedMessages
}

func IsMessage(
	i parser.ProtoInterface,
	name string,
) bool {

	for _, m := range i.Messages {
		if m.Name == name {
			return true
		}
	}
	return false
}

func AddImportToNestedFieldNames(
	fields []parser.Field,
	importName string,
) []parser.Field {

	retVal := make([]parser.Field, 0, len(fields))
	for _, f := range fields {
		if f.IsNestedMessage {
			f.Type = importName + "." + f.Type
		}
		retVal = append(retVal, f)
	}
	return retVal
}

func MethodUsesNestedMessages(
	i parser.ProtoInterface,
	methodName string,
) (bool, error) {

	reqFields, err := MethodRequestFields(i, methodName)
	if err != nil {
		return false, err
	}

	for _, f := range reqFields {
		if f.IsNestedMessage {
			return true, nil
		}
	}

	resFields, err := MethodResponseFields(i, methodName)
	if err != nil {
		return false, err
	}

	for _, f := range resFields {
		if f.IsNestedMessage {
			return true, nil
		}
	}

	return false, nil
}
