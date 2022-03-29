.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: build
build: build-migrate build-worker build-clock build-web

.PHONY: run
run:
	ENV=dev go run cmd/migrate/main.go
	ENV=dev go run cmd/worker/main.go &
	ENV=dev go run cmd/clock/main.go &
	ENV=dev go run cmd/web/main.go

.PHONY: build-migrate
build-migrate:
	go build -o dripfile-migrate cmd/migrate/main.go

.PHONY: run-migrate
run-migrate:
	ENV=dev go run cmd/migrate/main.go

.PHONY: build-worker
build-worker:
	go build -o dripfile-worker cmd/worker/main.go

.PHONY: run-worker
run-worker:
	ENV=dev go run cmd/worker/main.go

.PHONY: build-clock
build-clock:
	go build -o dripfile-clock cmd/clock/main.go

.PHONY: run-clock
run-clock:
	ENV=dev go run cmd/clock/main.go

.PHONY: build-web
build-web:
	go build -o dripfile-web cmd/web/main.go

.PHONY: run-web
run-web:
	ENV=dev go run cmd/web/main.go

.PHONY: test
test: run-migrate
	go test -short -count=1 -v ./...

.PHONY: test-ui
test-ui: run-migrate
	go test -count=1 -v ./...

.PHONY: race
race: run-migrate
	go test -short -race -count=1 ./...

.PHONY: race-ui
race-ui: run-migrate
	go test -race -count=1 ./...

.PHONY: cover
cover: run-migrate
	go test -short -coverprofile=c.out -coverpkg=./... -count=1 ./...
	go tool cover -html=c.out

.PHONY: cover-ui
cover-ui: run-migrate
	go test -coverprofile=c.out -coverpkg=./... -count=1 ./...
	go tool cover -html=c.out

.PHONY: format
format:
	go fmt ./...

.PHONY: clean
clean:
	rm -fr dripfile-* c.out dist/
