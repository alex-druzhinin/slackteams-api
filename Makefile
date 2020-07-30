M = $(shell printf "\033[34;1m▶\033[0m")
modules:
	go mod download

build: $(info $(M) Building project...)
	CGO_ENABLED=0  go build -a

server: $(info $(M) Starting development server...)
	env `cat ./env/.env | xargs` go run .

lint: $(info $(M) Running long lint from revision...)
	golangci-lint run --new-from-rev=ff1883d927ce4cd1e06242735efe613f3d919817

test: $(info $(M) Running all tests)
	go test ./...

.PHONY: build server lint modules test
