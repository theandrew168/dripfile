.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: build
build:
	go build -o dripfile main.go

.PHONY: run
run: build
	ENV=dev ./dripfile web

.PHONY: run-worker
run-worker: backend
	ENV=dev ./dripfile worker

.PHONY: run-scheduler
run-scheduler: backend
	ENV=dev ./dripfile scheduler

.PHONY: run-migrate
run-migrate: backend
	ENV=dev ./dripfile migrate

.PHONY: update
update:
	go get -u ./...
	go mod tidy

.PHONY: test
test: run-migrate
	go test -count=1 ./...

.PHONY: race
race: run-migrate
	go test -race -count=1 ./...

.PHONY: cover
cover: run-migrate
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
