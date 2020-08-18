FROM golang:1.15-alpine

RUN apk update && apk upgrade && \
    apk add build-base

ADD . /go/src/chat

WORKDIR /go/src/chat/chat/cmd/client

RUN go build

CMD ./client
