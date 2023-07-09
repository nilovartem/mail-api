build:
	go build -o ./bin/mailapi ./cmd/mailapi
run:
	./bin/mailapi
.DEFAULT_GOAL := build