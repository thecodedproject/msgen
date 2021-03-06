package parser

type ProtoInterface struct {

	ServiceName string
	ProtoPackage string
	Methods []Method
	Messages []Message
	//Enums []Enum
}

type Method struct {
	Name string
	RequestMessage string
	ResponseMessage string
}

type Message struct {
	Name string
	Fields []Field
}

type Field struct {
	Name string
	Type string
	IsNestedMessage bool
}
