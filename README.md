# Signature Service - Coding Challenge

## Instructions

### Local run

1. Run `go generate ./...` to generate the API boilerplate from the OpenAPI Spec
2. Run `go run main.go` to start the service

### Docker

1. Build with: `docker build --rm -t signing-service-challenge:0.0.1 .`
2. Run with: `docker run --rm -p 8080:8080 signing-service-challenge:0.0.1`

### Running tests

1. Run `go generate ./...` to generate the API boilerplate from the OpenAPI Spec
2. Run `docker run --pull always --rm -v ${PWD}:/local -u $(id -u) -w /local vektra/mockery` to generater mocks with mockery
3. Run `go test ./...` to run tests

## Configurations

Only parameter available is via the boolean env variable `CORS_ENABLED`, when true-ish (`strconv.ParseBool` way) Cross Origin Requests preflight check will be honored.

This is handy to lookup and try out APIs using the online [SwaggerUI](https://petstore3.swagger.io/?url=http://127.0.0.1:8080/api/v1/openapi.yaml).

Defaults to false.

## OpenAPI specification

The specification is available in [openapi/openapi.yaml](openapi/openapi.yaml) file or at URL [http://127.0.0.1:8080/api/v1/openapi.yaml](http://127.0.0.1:8080/api/v1/openapi.yaml) in a running application.

## Verification

The compliancy of the implementation could be verified by checking if the chaining of the responses is correctly respected (regardless of the response order).

1. Create a device
2. Validate its counter is 0 and lastSignature == base64(device.id)
// loop
3. Sign a transaction
4. Verify the transaction
5. Save the response
// end loop
6. Sort the responses by counter
7. Verify that counters are monotonically increasing and signatures are chained

see [domain/device_test.go:TestSign](domain/device_test.go#L136) for a test following above algorithm

## Considerations

* New signing algorithms can be added to the crypto package (implementing crypto/generation.go interfaces) and registered via init function
* A relational DB storage can be implemented using the Storage interface, the function mapping could be:

  * List -> SELECT ...
  * Get -> SELECT ... WHERE id = uuid
  * Add -> INSERT ...
  * Put -> UPDATE ...
