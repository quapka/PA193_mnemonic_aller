GO = go
BUILD_CMD = go build
INSTALL_CMD = go install

CMD_UTILITY = bip39

BIN_DIR = ./bin
SRC_DIR = ./src
CMD_DIR = ./cmd
PKG_DIR = ./pkg
EX_DIR  =  examples

EX1 = entropy-to-phrase-and-seed
EX2 = phrase-to-entropy-and-seed
EX3 = verify-phrase-and-seed

.PHONY: clean all

all: $(CMD_UTILITY)

$(CMD_UTILITY):
	$(BUILD_CMD) -o $(BIN_DIR)/$(CMD_UTILITY) $(CMD_DIR)/$@/*.go

install-pkg:
	$(INSTALL_CMD) ../PA193_mnemonic_aller/$(PKG_DIR)/mnemonic

test:
	$(GO) test -v $(PKG_DIR)/mnemonic/*.go

clean:
	rm -f $(BIN_DIR)/*

build-examples: $(EX1) $(EX2) $(EX3)

$(EX1):
	$(BUILD_CMD) -o $(BIN_DIR)/$@ $(EX_DIR)/$@/*.go

$(EX2):
	$(BUILD_CMD) -o $(BIN_DIR)/$@ $(EX_DIR)/$@/*.go

$(EX3):
	$(BUILD_CMD) -o $(BIN_DIR)/$@ $(EX_DIR)/$@/*.go
