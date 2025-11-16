FROM golang:1.24.3-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main cmd/http_server/main.go

ENTRYPOINT ["/app/main"]