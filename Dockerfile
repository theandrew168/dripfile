# syntax=docker/dockerfile:1

# Based on:
# https://docs.docker.com/language/golang/build-images/

## Build
FROM golang:1.20-buster AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -o /dripfile main.go

## Test
FROM build AS test
RUN go test -v -short ./...

## Deploy
FROM gcr.io/distroless/base-debian11
WORKDIR /
COPY --from=build /dripfile /dripfile
EXPOSE 5000
USER nonroot:nonroot
ENTRYPOINT ["/dripfile"]
