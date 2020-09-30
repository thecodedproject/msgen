* Nested proto types in method request/response types (create a new Go type for each nested type)

  * Type conversions to return errors

* Using nested types from other microservices

* Using custom go types (with custom conversion); e.g. a shopspring.Decimal

	* Define custom proto messages in another package:
    * assume that they will have *ToProto* *FromProto* funcs - parse that package to look for go funcs with that sig, and get go type that it is converted to (as well as it's import)
    * Store the go type and import (and conversion func/import) on the parser.Field
    * Nested messages can be handled in almost the same way (can't parse the go package because, at parse time, the conversion funcs wont exist yet... but everything else can work the same)
    * Also special case other types - e.g. built in types - float, `repeated`/slices and other custom types like google.timestamp

* Split out client_test file into seperate files (one for each client function) - not sure where the common code should go yet...
