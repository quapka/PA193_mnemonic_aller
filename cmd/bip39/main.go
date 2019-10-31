// Project: PA193_mnemonic_aller
// Maintainers UCO: 408788 497391 497577
// Description: Main file to test the API with command line

// Exemplar command line utility calculating the bip39 phrases and seeds.
// It can also verify whether phrase and seed match.
// Below you can see detailed help message. Functionality is chosen implicitly
// by observing the given flags.
//
//   $ ./bip39 --help
//   Usage of ./bip39:
//     -entropy string
//           Entropy to get phrase and seed (can't be set with --phrase)
//     -passphrase string
//           Passphrase to be used to generate the seed
//     -phrase string
//           Phrase to get entropy and seed (can't be set with --entropy)
//     -seed string
//           Seed to be provided with phrase to verify them (requires --phrase to be set)
//     -wordlist string
//           Path to wordlist (required)
//
package main

import (
	"flag"
	"fmt"
	"github.com/quapka/PA193_mnemonic_aller/pkg/mnemonic"
	"log"
	"os"
)

func displayMissingArg(msg string) {
	wrote, err := os.Stderr.WriteString(msg + "\nUsage:\n")
	if err != nil {
		log.Fatal(err)
	}
	for wrote != len(msg) {
		n, err := os.Stderr.WriteString(msg[wrote:])
		if err != nil {
			log.Fatal(err)
		}
		wrote += n
	}

	flag.PrintDefaults()
	os.Exit(1)
}

func main() {

	phrasePtr := flag.String("phrase", "",
		"Phrase to get entropy and seed (can't be set with --entropy)")
	entropyPtr := flag.String("entropy", "",
		"Entropy to get phrase and seed (can't be set with --phrase)")
	passphrasePtr := flag.String("passphrase", "",
		"Passphrase to be used to generate the seed")
	seedPtr := flag.String("seed", "",
		"Seed to be provided with phrase to verify them "+
			"(requires --phrase to be set)")
	wordlistFilePtr := flag.String("wordlist", "",
		"Path to wordlist (required)")

	flag.Parse()

	// Error handling
	if *wordlistFilePtr == "" {
		displayMissingArg("Need to set --wordlist")
	} else if *phrasePtr == "" && *entropyPtr == "" {
		displayMissingArg("Need to set --phrase or --entropy")
	} else if *phrasePtr != "" && *entropyPtr != "" {
		displayMissingArg("Can't set --phrase and --entropy")
	} else if *seedPtr != "" && *phrasePtr == "" {
		displayMissingArg("Requires to set --phrase")
	}

	if *seedPtr != "" {
		// Verify phrase and seed
		match, err := mnemonic.VerifyPhraseAndSeed(*phrasePtr,
			*passphrasePtr, *seedPtr)
		if err != nil {
			log.Fatal(err)
		}
		if !match {
			fmt.Println("The phrase and the seed do not match.")
		} else {
			fmt.Println("The phrase and the seed match each other.")
		}

	} else if *phrasePtr != "" {
		// Get the entropy and seed from the phrase
		if entropy, seed, err := mnemonic.PhraseToEntropyAndSeed(*phrasePtr,
			*passphrasePtr, *wordlistFilePtr); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("From phrase:", *phrasePtr)
			fmt.Println("Entropy:", entropy)
			fmt.Println("Seed:", seed)
		}

	} else {
		// Get the phrase and seed from the entropy
		if phrase, seed, err := mnemonic.EntropyToPhraseAndSeed(*entropyPtr,
			*passphrasePtr, *wordlistFilePtr); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("From entropy: ", *entropyPtr)
			fmt.Println("From passphrase: ", *passphrasePtr)
			fmt.Println("Phrase: ", phrase)
			fmt.Println("Seed: ", seed)
		}
	}
}
