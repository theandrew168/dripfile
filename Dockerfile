# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-buster AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /dripfile main.go

## Deploy
FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /dripfile /dripfile
EXPOSE 5000
USER nonroot:nonroot
ENTRYPOINT ["/dripfile"]
