.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: build
build: frontend backend

node_modules:
	npm install

.PHONY: frontend
frontend: node_modules
	npm run build

.PHONY: backend
backend: frontend
	go build -o dripfile main.go

.PHONY: watch-frontend
watch-frontend: node_modules
	npm run watch

.PHONY: run-backend
run-backend: frontend
	DEBUG=1 go run main.go

.PHONY: run
run: watch-frontend run-backend

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
format-frontend:
	npm run format

.PHONY: format-backend
format-backend:
	gofmt -l -s -w .

.PHONY: format
format: format-frontend format-backend

.PHONY: update
update:
	go get -u ./...
	go mod tidy

.PHONY: clean
clean:
	rm -fr dripfile main c.out node_modules/ public/index.js dist/
