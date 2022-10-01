FROM golang:1.18.1 AS base

WORKDIR /the-drink-almanac-api

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -mod=mod -o /api

ENTRYPOINT [ "/api" ]