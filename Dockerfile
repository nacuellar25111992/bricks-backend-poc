FROM golang:1.18-alpine as builder

ARG REVISION

RUN mkdir -p /bricks-backend-poc/

WORKDIR /bricks-backend-poc

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags "-s -w \
    -X github.com/nacuellar25111992/bricks-backend-poc/internal/version.REVISION=${REVISION}" \
    -a -o bin/bricks-backend-poc cmd/bricks-backend-poc/*