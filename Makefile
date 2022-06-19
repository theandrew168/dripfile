.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: build
build:
	go build -o dripfile main.go

.PHONY: web
web: build migrate
	ENV=dev ./dripfile web

.PHONY: worker
worker: build migrate
	ENV=dev ./dripfile worker

.PHONY: scheduler
scheduler: build migrate
	ENV=dev ./dripfile scheduler

.PHONY: migrate
migrate: build
	ENV=dev ./dripfile migrate

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
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run --fast --issues-exit-code 0

.PHONY: format
format:
	gofmt -l -s -w .

.PHONY: clean
clean:
	rm -fr dripfile c.out dist/
