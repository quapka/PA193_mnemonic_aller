// Project: PA193_mnemonic_aller
// Maintainers UCO: 408788, 497391, 497577
// Description: Mnemonic API

package mnemonic

import (
	"fmt"
	"strings"
)

const lowerENTBound = 128
const upperENTBound = 256

// entropy can be of various length, therefore it needs to a slice not an array
// FIXME check whether the underlying entropy array is changed, if so make a copy
// params:
// string entropy is a string of hexadecimal values
func EntropyToPhraseAndSeed(entropy string, dictFilepath string) (phrase, seed string, err error) {
	bytes, err := cleanInputEntropy(entropy)
	if err != nil {
		return "", "", err
	}
	// create binary string
	binary := convertToBinary(bytes)
	// checkSum
	checkSum := calculateCheckSum(bytes)
	binary += checkSum
	// create groups
	groups, _ := createGroups(binary)
	// create the indices
	indices, _ := createIndices(groups)
	// open the wordlist file
	wordList, err := loadWordlist(dictFilepath)
	if err != nil {
		return "", "", err
	}
	phraseWords, _ := createPhraseWords(indices, wordList)

	seed = ""

	return strings.Join(phraseWords, " "), seed, nil
}

func PhraseToEntropyAndSeed() {
	// TODO
	fmt.Println("TODO entropy_to_mnemonic")
}

func VerifyPhraseAndSeed() {
	// TODO
	fmt.Println("TODO entropy_to_mnemonic")
}
