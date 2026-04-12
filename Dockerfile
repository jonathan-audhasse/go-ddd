FROM golang:1.26.2-alpine

RUN apk add zsh git postgresql-client curl make

WORKDIR /go-ddd

COPY . /go-ddd

# install mockery to generate mock file
RUN cd goddd && go install github.com/vektra/mockery/v3@v3.7.0 github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0
RUN cd goddd && go mod tidy
#run server
#CMD go run goddd/cmd/api/main.go
