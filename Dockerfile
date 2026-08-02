FROM golang:1.26.2-alpine AS builder

WORKDIR /src

COPY goddd .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/migrate ./cmd/migrate


FROM alpine:3.22

RUN apk add --no-cache ca-certificates

COPY --from=builder /out/api /app/api
COPY --from=builder /out/migrate /app/migrate

WORKDIR /app