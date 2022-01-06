.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: css
css:
	tailwindcss -m -i static/css/tailwind.input.css -o static/css/tailwind.min.css

.PHONY: build
build: build-web build-worker build-scheduler

.PHONY: build-web
build-web: css
	go build -o dripfile-web cmd/web/main.go

.PHONY: run-web
run-web:
	ENV=dev go run cmd/web/main.go -conf internal/test/dripfile.conf &
	tailwindcss --watch -m -i static/css/tailwind.input.css -o static/css/tailwind.min.css

.PHONY: build-worker
build-worker:
	go build -o dripfile-worker cmd/worker/main.go

.PHONY: run-worker
run-worker:
	ENV=dev go run cmd/worker/main.go -conf internal/test/dripfile.conf

.PHONY: build-scheduler
build-scheduler:
	go build -o dripfile-scheduler cmd/scheduler/main.go

.PHONY: run-scheduler
run-scheduler:
	ENV=dev go run cmd/scheduler/main.go -conf internal/test/dripfile.conf

.PHONY: test
test:
	go run main.go -conf internal/test/dripfile.conf -migrate
	go test -count=1 -v ./...

.PHONY: race
race:
	go run main.go -conf internal/test/dripfile.conf -migrate
	go test -race -count=1 ./...

.PHONY: cover
cover:
	go run main.go -conf internal/test/dripfile.conf -migrate
	go test -coverprofile=c.out -coverpkg=./... -count=1 ./...
	go tool cover -html=c.out

.PHONY: format
format:
	go fmt ./...

.PHONY: clean
clean:
	rm -fr dripfile-* c.out dist/
