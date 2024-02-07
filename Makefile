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

.PHONY: run-frontend
run-frontend: frontend
	npm run dev

.PHONY: run-backend
run-backend: frontend
	DEBUG=1 go run main.go

# run the backend and frontend concurrently (requires at least "-j2") 
.PHONY: run
run: run-frontend run-backend

.PHONY: build-image
build-image:
	docker build -t dripfile .

.PHONY: run-image
run-image:
	docker run -p 5000:5000 dripfile

.PHONY: migrate
migrate:
	go run main.go -migrate

.PHONY: test
test: migrate
	go test -count=1 ./...

.PHONY: race
race:
	go test -race -count=1 ./...

.PHONY: cover
cover:
	go test -coverprofile=c.out -coverpkg=./... -count=1 ./...
	go tool cover -html=c.out

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
	npm update --save

.PHONY: update-backend
update-backend:
	go get -u ./...
	go mod tidy

.PHONY: update
update: update-frontend update-backend

.PHONY: clean
clean:
	rm -fr dripfile main c.out public/index.js public/index.css dist/
