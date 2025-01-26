# go-ddd
Domain Driven Design pattern in go

## Requirements

The following tools need to be installed in order to set up your local environment:
- [docker](https://docs.docker.com/get-started/get-docker)
- [go](https://go.dev/doc/install) (optional) if you do not develop in the container

## Mount your server

Local usage:
```sh
# to mount your service
make test
```
```sh
# to shut down your service
make test
```

## Tests

To run all tests (unit, integration and e2e tests) run the following command: 
```sh
make test
```

! you need to have [mockery](https://vektra.github.io/mockery/latest) installed locally in order to update the code.

## Help

For more help run:
```sh
make help
```
