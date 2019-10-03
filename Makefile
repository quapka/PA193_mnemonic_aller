
BUILD_CMD=go build

BIN_NAME=aller.out

BIN_DIR=./bin
SRC_DIR=./src

export GOPATH=$(shell pwd)

build:
	$(BUILD_CMD) -o $(BIN_DIR)/$(BIN_NAME) $(SRC_DIR)/*.go

start: build
	./bin/$(BIN_NAME)

test:
	go test -cover mnemonic


clean:
	rm -f bin/*

