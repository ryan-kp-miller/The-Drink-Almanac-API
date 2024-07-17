FROM golang:1.18.1 AS base

WORKDIR /the-drink-almanac-api

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -mod=mod -o /api

FROM scratch
COPY --from=base /api /usr/bin/api
ENTRYPOINT [ "api" ]