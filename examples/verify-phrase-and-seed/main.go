// Minimal working example of using mnemonic.VerifyPhraseAndSeed
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
	// set the passphrase
	passphrase := "TREZOR"
	// set the seed
	seed := "c55257c360c07c72029aebc1b53c05ed0362ada38" +
		"ead3e3e9efa3708e53495531f09a6987599d18264" +
		"c1e1c92f2cf141630c7a3c4ab7c81b2f001698e7463b04"

	match, err := mnemonic.VerifyPhraseAndSeed(phrase, passphrase, seed)
	if err != nil {
		log.Fatal(err)
	}
	if !match {
		fmt.Println("The phrase and seed do not match.")
	}
	fmt.Println("The phrase and seed do match.")
}
