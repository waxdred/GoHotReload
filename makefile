
BINARY_NAME := gohot
SOURCES := $(shell find . -name '*.go')

.DEFAULT_GOAL := all

all: build exec

build: $(BINARY_NAME)

$(BINARY_NAME): $(SOURCES)
	go build -o $(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)

run:
	go run .

exec: build
	./$(BINARY_NAME)

re: clean build exec

.PHONY: all build clean run
