FROM golang:1.23-alpine

RUN apk add zsh git postgresql-client curl make

WORKDIR /go-ddd

COPY . /go-ddd

# install mockery to generate mock file
RUN go install github.com/vektra/mockery/v2@v2.51.1
RUN go mod tidy
#run server
#CMD go run src/main.go
