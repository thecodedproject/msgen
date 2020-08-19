package parser

import(
	"log"
	"github.com/emicklei/proto"
	"os"
)

func Parse(filePath string) (ProtoInterface, error) {

	reader, err := os.Open(filePath)
	if err != nil {
		return ProtoInterface{}, err
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		return ProtoInterface{}, err
	}

	var protoInterface ProtoInterface

	proto.Walk(
		definition,
		proto.WithEnum(enumHandler(&protoInterface)),
		proto.WithImport(importHandler(&protoInterface)),
		proto.WithMessage(messageHandler(&protoInterface)),
		proto.WithOneof(oneofHandler(&protoInterface)),
		proto.WithOption(optionHandler(&protoInterface)),
		proto.WithPackage(packageHandler(&protoInterface)),
		proto.WithRPC(rpcHandler(&protoInterface)),
		proto.WithService(serviceHandler(&protoInterface)),
	)

	return protoInterface, nil
}

func enumHandler(ms *ProtoInterface) func(*proto.Enum) {

	return func(*proto.Enum) {
	}
}

func importHandler(ms *ProtoInterface) func(*proto.Import) {

	return func(*proto.Import) {
	}
}

func messageHandler(ms *ProtoInterface) func(*proto.Message) {

	return func(m *proto.Message) {

		message := Message{
			Name: m.Name,
		}

		for _, f := range m.Elements {
			field, ok := f.(*proto.NormalField)
			if !ok {
				log.Fatal("Only normal fields supported")
			}

			message.Fields = append(message.Fields, Field{
				Name: field.Name,
				Type: field.Type,
			})
		}

		ms.Messages = append(ms.Messages, message)
	}
}

func oneofHandler(ms *ProtoInterface) func(*proto.Oneof) {

	return func(*proto.Oneof) {

		log.Fatal("OneOf not supported")
	}
}

func optionHandler(ms *ProtoInterface) func(*proto.Option) {

	return func(*proto.Option) {
	}
}

func packageHandler(ms *ProtoInterface) func(*proto.Package) {

	return func(p *proto.Package) {

		ms.ProtoPackage = p.Name
	}
}

func rpcHandler(ms *ProtoInterface) func(*proto.RPC) {

	return func(r *proto.RPC) {

		ms.Methods = append(ms.Methods, Method{
			Name: r.Name,
			RequestMessage: r.RequestType,
			ResponseMessage: r.ReturnsType,
		})
	}
}

func serviceHandler(ms *ProtoInterface) func(*proto.Service) {

	return func(*proto.Service) {
	}
}

