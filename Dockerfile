# syntax=docker/dockerfile:1

## Build
FROM golang:1.18-buster AS build
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o /dripfile main.go

## Deploy
FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /dripfile /dripfile
EXPOSE 5000
USER nonroot:nonroot
ENTRYPOINT ["/dripfile"]
