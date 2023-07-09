.PHONY: test
build:
	go build -o ./bin/mailapi ./cmd/mailapi
run:
	./bin/mailapi
test:
	go test ./internal/app/model
	go test ./internal/app/mailapi
all: test build run
.DEFAULT_GOAL := all