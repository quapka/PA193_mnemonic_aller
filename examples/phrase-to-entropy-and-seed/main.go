package main

import (
	"fmt"
	"github.com/quapka/PA193_mnemonic_aller/pkg/mnemonic"
	"log"
)

func main() {
	// example values taken from
	// https://github.com/trezor/python-mnemonic/blob/master/vectors.json
	// set the phrase
	phrase := "abandon abandon abandon abandon abandon abandon " +
		"abandon abandon abandon abandon abandon about"
	// optional set the passphrase, can be empty ""
	passphrase := "TREZOR"
	// set the filepath to the wordlist that should be used
	wordlistFilepath := "../../wordlists/english.txt"
	// calculate the entropy and seed
	entropy, seed, err := mnemonic.PhraseToEntropyAndSeed(phrase,
		passphrase, wordlistFilepath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Entropy:\n%s\n", entropy)
	fmt.Printf("Seed:\n%s\n", seed)
}
