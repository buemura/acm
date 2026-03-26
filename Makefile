APP_NAME := acm

.PHONY: build run clean help

build:
	go build -o $(APP_NAME) .

run: build
	./$(APP_NAME)

clean:
	rm -f $(APP_NAME)

help: build
	./$(APP_NAME) help
