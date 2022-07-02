.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: build
build:
	go build -o dripfile main.go

.PHONY: web
web: migrate
	ENV=dev go run main.go web

.PHONY: worker
worker: migrate
	ENV=dev go run main.go worker

.PHONY: scheduler
scheduler: migrate
	ENV=dev go run main.go scheduler

.PHONY: migrate
migrate:
	ENV=dev go run main.go migrate

.PHONY: update
update:
	go get -u ./...
	go mod tidy

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
release:
	goreleaser release --snapshot --rm-dist

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run --fast

.PHONY: format
format:
	gofmt -l -s -w .

.PHONY: clean
clean:
	rm -fr dripfile c.out dist/
