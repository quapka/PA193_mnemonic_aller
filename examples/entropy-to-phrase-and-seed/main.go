// Minimal working example of using mnemonic.EntropyToPhraseAndSeed
package main

import (
	"fmt"
	"github.com/quapka/PA193_mnemonic_aller/pkg/mnemonic"
	"log"
)

func main() {
	// example values taken from
	// https://github.com/trezor/python-mnemonic/blob/master/vectors.json
	// set the entropy
	entropy := "00000000000000000000000000000000"
	// optional: set the passphrase, can be empty ""
	passphrase := "TREZOR"
	// set the filepath to the wordlist that should be used
	wordlistFilepath := "../../wordlists/english.txt"
	// calculate the phrase and seed
	phrase, seed, err := mnemonic.EntropyToPhraseAndSeed(entropy,
		passphrase, wordlistFilepath)
	if err != nil {
		// exit in case there was an error
		log.Fatal(err)
	}
	fmt.Printf("Phrase:\n%s\n", phrase)
	fmt.Printf("Seed:\n%s\n", seed)
}
