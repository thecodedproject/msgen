package files

import(
	"github.com/thecodedproject/msgen/generator/files/api"
	"github.com/thecodedproject/msgen/generator/files/client_grpc"
	"github.com/thecodedproject/msgen/generator/files/client_logical"
	"github.com/thecodedproject/msgen/generator/files/client_test_file"
	"github.com/thecodedproject/msgen/generator/files/ops_backends"
	"github.com/thecodedproject/msgen/generator/files/ops_functions"
	"github.com/thecodedproject/msgen/generator/files/rpc_server"
	"github.com/thecodedproject/msgen/parser"
	//"path"
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
		return err
	}

	err = client_grpc.Generate(
		serviceRootImportPath,
		i,
		outputDir,
	)
	if err != nil {
		return err
	}

	err = client_logical.Generate(
		serviceRootImportPath,
		i,
		outputDir,
	)
	if err != nil {
		return err
	}

	err = client_test_file.Generate(
		serviceRootImportPath,
		i,
		outputDir,
	)
	if err != nil {
		return err
	}

	err = ops_backends.Generate(
		serviceRootImportPath,
		outputDir,
	)
	if err != nil {
		return err
	}

	err = ops_functions.Generate(
		serviceRootImportPath,
		i,
		outputDir,
	)
	if err != nil {
		return err
	}

	err = rpc_server.Generate(
		serviceRootImportPath,
		i,
		outputDir,
	)
	if err != nil {
		return err
	}

	return nil
}
