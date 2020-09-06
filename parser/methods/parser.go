package methods

import(
	"github.com/emicklei/proto"
	msgen_parser "github.com/thecodedproject/msgen/parser"
	"github.com/thecodedproject/msgen/parser/proto_helpers"
	"log"
	"os"
)

func Parse(filePath string) (msgen_parser.ProtoInterface, error) {

	reader, err := os.Open(filePath)
	if err != nil {
		return msgen_parser.ProtoInterface{}, err
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		return msgen_parser.ProtoInterface{}, err
	}

	var protoInterface msgen_parser.ProtoInterface

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

	protoInterface = markNestedMessages(protoInterface)

	return protoInterface, nil
}

func enumHandler(ms *msgen_parser.ProtoInterface) func(*proto.Enum) {

	return func(*proto.Enum) {
	}
}

func importHandler(ms *msgen_parser.ProtoInterface) func(*proto.Import) {

	return func(*proto.Import) {
	}
}

func messageHandler(ms *msgen_parser.ProtoInterface) func(*proto.Message) {

	return func(m *proto.Message) {

		message := msgen_parser.Message{
			Name: m.Name,
		}

		for _, f := range m.Elements {
			field, ok := f.(*proto.NormalField)
			if !ok {
				log.Fatal("Only normal fields supported")
			}

			message.Fields = append(message.Fields, msgen_parser.Field{
				Name: field.Name,
				Type: field.Type,
			})
		}

		ms.Messages = append(ms.Messages, message)
	}
}

func oneofHandler(ms *msgen_parser.ProtoInterface) func(*proto.Oneof) {

	return func(*proto.Oneof) {

		log.Fatal("OneOf not supported")
	}
}

func optionHandler(ms *msgen_parser.ProtoInterface) func(*proto.Option) {

	return func(*proto.Option) {
	}
}

func packageHandler(ms *msgen_parser.ProtoInterface) func(*proto.Package) {

	return func(p *proto.Package) {

		ms.ProtoPackage = p.Name
	}
}

func rpcHandler(ms *msgen_parser.ProtoInterface) func(*proto.RPC) {

	return func(r *proto.RPC) {

		ms.Methods = append(ms.Methods, msgen_parser.Method{
			Name: r.Name,
			RequestMessage: r.RequestType,
			ResponseMessage: r.ReturnsType,
		})
	}
}

func serviceHandler(ms *msgen_parser.ProtoInterface) func(*proto.Service) {

	return func(s *proto.Service) {
		ms.ServiceName = s.Name
	}
}

func markNestedMessages(i msgen_parser.ProtoInterface) msgen_parser.ProtoInterface {

	for iM := range i.Messages {
		for iF := range i.Messages[iM].Fields {
			if proto_helpers.IsMessage(i, i.Messages[iM].Fields[iF].Type) {
				i.Messages[iM].Fields[iF].IsNestedMessage = true
			}
		}
	}
	return i
}
