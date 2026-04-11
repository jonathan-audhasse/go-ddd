FROM golang:1.26.2-alpine

RUN apk add zsh git postgresql-client curl make

WORKDIR /go-ddd

COPY . /go-ddd

# install mockery to generate mock file
RUN cd src && go install github.com/vektra/mockery/v3@v3.7.0
RUN cd src && go mod tidy
#run server
#CMD go run src/main.go
