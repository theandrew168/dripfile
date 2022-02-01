.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: build
build: build-migrate build-web build-worker build-scheduler build-tap

.PHONY: run
run:
	ENV=dev go run cmd/migrate/main.go
	ENV=dev go run cmd/web/main.go &
	ENV=dev go run cmd/worker/main.go &
	ENV=dev go run cmd/scheduler/main.go &
	tailwindcss --watch -m -i tailwind.input.css -o internal/static/static/css/tailwind.min.css

.PHONY: build-migrate
build-migrate:
	go build -o dripfile-migrate cmd/migrate/main.go

.PHONY: run-migrate
run-migrate:
	ENV=dev go run cmd/migrate/main.go

.PHONY: build-css
build-css:
	tailwindcss -m -i tailwind.input.css -o internal/static/static/css/tailwind.min.css

.PHONY: build-web
build-web: build-css
	go build -o dripfile-web cmd/web/main.go

.PHONY: run-web
run-web:
	ENV=dev go run cmd/web/main.go &
	tailwindcss --watch -m -i tailwind.input.css -o internal/static/static/css/tailwind.min.css

.PHONY: build-worker
build-worker:
	go build -o dripfile-worker cmd/worker/main.go

.PHONY: run-worker
run-worker:
	ENV=dev go run cmd/worker/main.go

.PHONY: build-scheduler
build-scheduler:
	go build -o dripfile-scheduler cmd/scheduler/main.go

.PHONY: run-scheduler
run-scheduler:
	ENV=dev go run cmd/scheduler/main.go

.PHONY: build-tap
build-tap:
	go build -o dripfile-tap cmd/tap/main.go

.PHONY: run-tap
run-tap:
	ENV=dev go run cmd/tap/main.go

.PHONY: test
test: run-migrate
	go test -short -count=1 -v ./...

.PHONY: test-ui
test-ui: run-migrate
	go test -count=1 -v ./...

.PHONY: race
race: run-migrate
	go test -race -count=1 ./...

.PHONY: cover
cover: run-migrate
	go test -coverprofile=c.out -coverpkg=./... -count=1 ./...
	go tool cover -html=c.out

.PHONY: format
format:
	go fmt ./...

.PHONY: clean
clean:
	rm -fr dripfile-* c.out dist/
