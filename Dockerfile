FROM golang:1.22-alpine

RUN apk add zsh git postgresql-client curl

WORKDIR /go-ddd

COPY . /go-ddd

RUN go mod tidy
#run server
#CMD go run src/main.go
