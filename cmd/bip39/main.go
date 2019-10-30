// Project: PA193_mnemonic_aller
// Maintainers UCO: 408788 497391 497577
// Description: Main file to test the API with command line

package main

import (
	"flag"
	"fmt"
	"mnemonic"
	"os"
)

func displayMissingArg(msg string) {
	os.Stderr.WriteString(msg + "\nUsage:\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {

	phrasePtr := flag.String("phrase", "", "Phrase to get entropy and seed (can't be set with --entropy)")
	entropyPtr := flag.String("entropy", "", "Entropy to get phrase and seed (can't be set with --phrase)")
	passphrasePtr := flag.String("passphrase", "", "Passphrase to get phrase and seed")
	seedPtr := flag.String("seed", "", "Seed to be provided with phrase to verify them (requires --phrase to be set)")
	wordlistFilePtr := flag.String("wordlist", "", "Path to wordlist (required)")

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
		mnemonic.VerifyPhraseAndSeed("TMP CHARACTER CHAIN", "NEED TO BE CHANGED", *wordlistFilePtr)

	} else if *phrasePtr != "" {
		// Get the entropy and seed from the phrase
		if entropy, e := mnemonic.PhraseToEntropyAndSeed(*phrasePtr, *wordlistFilePtr); e != nil {
			fmt.Println(e)
		} else {
			fmt.Println("From phrase:", *phrasePtr)
			fmt.Println("Entropy:", entropy)
			fmt.Println("Seed: TODO")
		}

	} else {
		// Get the phrase and seed from the entropy
		if phrase, seed, err := mnemonic.EntropyToPhraseAndSeed(*entropyPtr, *passphrasePtr, *wordlistFilePtr); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("From entropy: ", *entropyPtr)
			fmt.Println("From passphrase: ", *passphrasePtr)
			fmt.Println("Phrase: ", phrase)
			fmt.Println("Seed: ", seed)
		}
	}
}