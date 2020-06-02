FROM golang:latest as go
FROM golangci/golangci-lint:v1.27-alpine

RUN apk update && apk add --no-cache --update python3 && apk add make

COPY . /go/src/github.com/kichyr/cubeCliClient
WORKDIR /go/src/github.com/kichyr/cubeCliClient

RUN pip3 install -r ./test/requirements.txt

RUN make test-local
RUN make lint