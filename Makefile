.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: build
build: frontend backend

node_modules:
	npm install

.PHONY: frontend-types
frontend-types: node_modules
	npm run build-types

.PHONY: frontend-js
frontend-js: node_modules
	npm run build-js

.PHONY: frontend-css
frontend-css: node_modules
	npm run build-css

.PHONY: frontend
frontend: frontend-types frontend-js frontend-css

.PHONY: backend
backend: frontend
	go build -o dripfile main.go

.PHONY: run-frontend-js
run-frontend-js: node_modules
	npm run run-js

.PHONY: run-frontend-css
run-frontend-css: node_modules
	npm run run-css

.PHONY: run-frontend
run-frontend: run-frontend-js run-frontend-css

.PHONY: run-backend
run-backend:
	DEBUG=1 go run github.com/cosmtrek/air@latest

.PHONY: run
run: run-frontend run-backend

.PHONY: migrate
migrate:
	go run main.go -migrate

.PHONY: test
test: migrate
	go test -count=1 ./...

.PHONY: race
race: migrate
	go test -race -count=1 ./...

.PHONY: cover
cover: migrate
	go test -coverprofile=c.out -coverpkg=./... -count=1 ./...
	go tool cover -html=c.out

.PHONY: release
release: frontend
	goreleaser release --snapshot --clean

.PHONY: format-frontend
format-frontend: node_modules
	npm run format

.PHONY: format-backend
format-backend:
	gofmt -l -s -w .

.PHONY: format
format: format-frontend format-backend

.PHONY: update-frontend
update-frontend:
	npm update

.PHONY: update-backend
update-backend:
	go get -u ./...
	go mod tidy

.PHONY: update
update: update-frontend update-backend

.PHONY: clean
clean:
	rm -fr dripfile main c.out public/index.js public/index.css dist/
