ARG GOVERSION=1.15

FROM golang:${GOVERSION}-alpine

EXPOSE 9090

WORKDIR /go/src/github.com/pagarme/marshals/labs/vicki/desafioGo

RUN apk add --no-cache graphviz git gcc musl-dev
