FROM golang:1.13.4-alpine3.10 as builder
# gopath is /go
# workdir is /go
WORKDIR /go

# Workaround to speed up building with cache
RUN apk add --no-cache git &&\
    mkdir -p /go/src/github.com/kamilkoduo/digicart
WORKDIR /go/src/github.com/kamilkoduo/digicart
COPY ./go.mod ./go.mod
RUN go mod tidy

#Improvements
ENV CGO_ENABLED 0

COPY ./src /go/src/github.com/kamilkoduo/digicart/src

#RUN go test /go/src/github.com/kamilkoduo/digicart/src/rest
# suggested scenario: run container with code with go test command in a docker-compose with redis,
# then build binary and put to alpine

RUN go build -o /app/cart/build /go/src/github.com/kamilkoduo/digicart/src/main

FROM alpine:3.10.3

COPY --from=builder /app /app
