FROM golang:1.16.0-alpine3.13

RUN apk add build-base

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go mod download

RUN go build ./cmd/api

RUN go test ./...

EXPOSE 8080
CMD ["/app/api"]