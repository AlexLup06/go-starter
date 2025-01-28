ARG ALPINE_VERSION=3.19
ARG GO_VERSION=1.22.2

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as base

RUN apk update && apk add --no-cache postgresql git gcc libc-dev
ENV CGO_ENABLED=0
RUN go install github.com/rubenv/sql-migrate/...@latest

WORKDIR /migrations

COPY . .


CMD ["sql-migrate", "up", "-env=development" , "-config=dbconfig.yaml"]

