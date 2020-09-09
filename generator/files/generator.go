package files

import(
	"github.com/pkg/errors"
	"github.com/thecodedproject/msgen/generator/files/api"
	"github.com/thecodedproject/msgen/generator/files/client_grpc"
	"github.com/thecodedproject/msgen/generator/files/client_logical"
	"github.com/thecodedproject/msgen/generator/files/client_test_file"
	"github.com/thecodedproject/msgen/generator/files/ops_functions"
	"github.com/thecodedproject/msgen/generator/files/proto_type_conversion"
	"github.com/thecodedproject/msgen/generator/files/rpc_server"
	"github.com/thecodedproject/msgen/generator/files/state"
	"github.com/thecodedproject/msgen/generator/files/types"
	"github.com/thecodedproject/msgen/parser"
)

func Generate(
	serviceRootImportPath string,
	i parser.ProtoInterface,
	outputDir string,
) error {

	err := api.Generate(
		serviceRootImportPath,
		i,
		outputDir,
	)
	if err != nil {
		return errors.Wrap(err, "api.Generate")
	}

	err = client_grpc.Generate(
		serviceRootImportPath,
		i,
		outputDir,
	)
	if err != nil {
		return errors.Wrap(err, "client_grpc.Generate")
	}

	err = client_logical.Generate(
		serviceRootImportPath,
		i,
		outputDir,
	)
	if err != nil {
		return errors.Wrap(err, "client_logical.Generate")
	}

	err = client_test_file.Generate(
		serviceRootImportPath,
		i,
		outputDir,
	)
	if err != nil {
		return errors.Wrap(err, "client_test_file.Generate")
	}

	err = ops_functions.Generate(
		serviceRootImportPath,
		i,
		outputDir,
	)
	if err != nil {
		return errors.Wrap(err, "ops_functions.Generate")
	}

	err = proto_type_conversion.Generate(
		serviceRootImportPath,
		i,
		outputDir,
	)
	if err != nil {
		return errors.Wrap(err, "proto_type_conversion.Generate")
	}

	err = rpc_server.Generate(
		serviceRootImportPath,
		i,
		outputDir,
	)
	if err != nil {
		return errors.Wrap(err, "rpc_server.Generate")
	}

	err = state.Generate(
		serviceRootImportPath,
		outputDir,
	)
	if err != nil {
		return errors.Wrap(err, "state.Generate")
	}

	err = types.Generate(
		serviceRootImportPath,
		i,
		outputDir,
	)
	if err != nil {
		return errors.Wrap(err, "types.Generate")
	}

	return nil
}
