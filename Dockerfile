# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /rest-api-framework
ENV PROJECT_PATH='/rest-api-framework'
ENV ENV='docker'

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN apk add --update make
RUN make build_all

EXPOSE 8080

CMD [ "/rest-api-framework/bin/run-api" ]