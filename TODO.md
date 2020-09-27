* Nested proto types in method request/response types (create a new Go type for each nested type)

  * Make sure that the types import is not included if a nested type is not used

* Using nested types from other microservices

* Using custom go types (with custom conversion); e.g. a shopspring.Decimal

* Split out client_test file into seperate files (one for each client function) - not sure where the common code should go yet...
