# syntax=docker/dockerfile:1

FROM node:20 AS build-frontend
WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . ./
RUN npm run build-types
RUN npm run build-js
RUN npm run build-css


FROM golang:1.21 AS build-backend
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
COPY --from=build-frontend /app/public ./public/
RUN CGO_ENABLED=0 GOOS=linux go build -o /dripfile


FROM gcr.io/distroless/base-debian12
WORKDIR /

COPY --from=build-backend /dripfile /dripfile
COPY dripfile.docker.conf /dripfile.conf

USER nonroot:nonroot
EXPOSE 5000
ENTRYPOINT ["/dripfile"]
