
BUILD_CMD=go build

BIN_DIR=./bin
SRC_DIR=./src


all: main api

main:
	$(BUILD_CMD) -o $(BIN_DIR)/main.out $(SRC_DIR)/main.go

api:
	$(BUILD_CMD) -o $(BIN_DIR)/mnemonic.out $(SRC_DIR)/mnemonic.go

clean:
	rm -f *.out

