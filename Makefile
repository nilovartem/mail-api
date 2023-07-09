.PHONY: test
build:
	go build -o ./bin/mailapi ./cmd/mailapi
run:
	./bin/mailapi
test:
	go test ./internal/app/model
	go test ./internal/app/mailapi
all: test build run
#TODO: добавить скрипт для тестирования с curl
#curl -X POST --user "karlenko.anton@wb.ru:karlenko" localhost:8080/karlenko.anton@wb.ru
.DEFAULT_GOAL := build