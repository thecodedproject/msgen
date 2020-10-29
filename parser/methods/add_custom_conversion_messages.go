package methods

import(

)

// Add the appropriate field objects to the proto.Interface, containing:
//
// * That this is a custom conversion type (?)
// * Conversion function names (We must store them rather than infer them later - some types (e.g. google.Timestamp) will not follow the `ToProto`/`FromProto` standards)
// * Conversion functions import
// * The go type for this field
// * The import for the go type
//
// Must also check for:
//
// **IF** this is a field with custom convertors:
//
// * That the conversion functions exist
// * The that Go type exists (i.e. `protoc` has been run)
// * That the conversion functions both convert to/from the same go type
// * That the conversion functions follow the desired interface


func AddCustomConversionFields(
	packageDir string,
	fieldNames []string,
) error {

	return nil
}
