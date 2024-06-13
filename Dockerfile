FROM golang:1.22-alpine

RUN apk add zsh git make

WORKDIR /go-ddd

COPY . /go-ddd

#run server
#CMD go run src/main.go
