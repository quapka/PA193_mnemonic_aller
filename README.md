[![Go report](https://goreportcard.com/badge/github.com/quapka/PA193_mnemonic_aller)](https://goreportcard.com/report/github.com/quapka/PA193_mnemonic_aller)
[![Build Status](https://travis-ci.com/quapka/PA193_mnemonic_aller.svg?branch=master)](https://travis-ci.com/quapka/PA193_mnemonic_aller)
[![Coverage Status](https://coveralls.io/repos/github/quapka/PA193_mnemonic_aller/badge.svg?branch=master)](https://coveralls.io/github/quapka/PA193_mnemonic_aller?branch=master)

PA193_mnemonic_aller
====================

A golang implementation of the [BIP-39 specification](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki).

## Install

In order to download the package you need to run:

```
$ go get github.com/quapka/PA193_mnemonic_aller/pkg/mnemonic
```

## Make

The project wraps some `go build` and `go install` calls in `make` calls.


## Use the command line interface

Along with the package there is a simple command line interface, that allows you to use the full functionality of the library. However, this utility is not intended for production, rather serves as an example.

```
$ ./bip39 --help
Usage:
  -entropy string
        Entropy to get phrase and seed (can't be set with --phrase)
  -passphrase string
        Passphrase to get phrase and seed
  -phrase string
        Phrase to get entropy and seed (can't be set with --entropy)
  -seed string
        Seed to be provided with phrase to verify them (requires --phrase to be set)
  -wordlist string
        Path to wordlist (required)

$ ./bip39 --entropy "$(python -c "print '00' * 16")" --passphrase "TREZOR" --wordlist ../wordlists/english.txt
From entropy:  00000000000000000000000000000000
From passphrase:  TREZOR
Phrase:  abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about
Seed:  c55257c360c07c72029aebc1b53c05ed0362ada38ead3e3e9efa3708e53495531f09a6987599d18264c1e1c92f2cf141630c7a3c4ab7c81b2f001698e7463b04
```

In order to build the utility run:
```
# the executable will end up in the ./bin folder
$ make bip39
# or create the binary locally
$ cd cmd/bip39 && go build .
```

## Application Programming Interface

The API consists of simply three functions. You can see the signatures here. In order to

```
func EntropyToPhraseAndSeed(entropy, passphrase, dictFilepath string) (phrase, seed string, err error)

func PhraseToEntropyAndSeed(phrase, passphrase, wlfile string) (string, string, error)

func VerifyPhraseAndSeed(phrase, passphrase, seed string) (bool, error)
```

## Examples

There are three examples located in `examples/` directory. Each subfolder contains example code using one of the functions of the API. You can either build the examples separately by running `go build .` inside the directory e.g.:
```
$ cd examples/entropy-to-phrase-and-seed
$ go build .
```

Or build them using `make`:
```
# build all examples at once
$ make build-examples
# build only one exapmle
$ make verify-phrase-and-seed
```

## Credits

The test vectors come from the Python implementation made by the Trezor team: [https://github.com/trezor/python-mnemonic/blob/master/vectors.json](https://github.com/trezor/python-mnemonic/blob/master/vectors.json)

The English wordlist (in `worldlists/english.txt`)  is from the [BIP-39 spec](https://github.com/bitcoin/bips/blob/master/bip-0039/bip-0039-wordlists.md).

