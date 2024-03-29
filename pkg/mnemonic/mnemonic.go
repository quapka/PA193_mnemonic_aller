// Project: PA193_mnemonic_aller
// Maintainers UCO: 408788, 497391, 497577
// Description: Mnemonic API

// Package mnemonic provides three exported functions for creating and
// verifying BIP39 phrases and seeds according to
// https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki.
package mnemonic

import (
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"strings"
)

// EntropyToPhraseAndSeed accepts entropy, optional passphrase and path to the
// word file. It then creates the phrase and the seed.
func EntropyToPhraseAndSeed(entropy, passphrase,
	wlFile string) (phrase, seed string, err error) {
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
	wordList, err := loadWordlist(wlFile)
	if err != nil {
		return "", "", err
	}
	// create the mnemonic phrase
	phraseWords, _ := createPhraseWords(indices, wordList)
	phrase = strings.Join(phraseWords, " ")
	// create the seed value
	seedBytes, _ := phraseToSeed(phrase, passphrase)
	seed = hex.EncodeToString(seedBytes)

	return phrase, seed, nil
}

// PhraseToEntropyAndSeed take the phrase, optional passphrase and path to the
// word file. It then creates the entropy and the seed.
func PhraseToEntropyAndSeed(phrase, passphrase,
	wlFile string) (string, string, error) {

	var wBytes [2]byte

	// Use a big int for arbitrary-precision arithmetic of big numbers
	bg := big.NewInt(0)
	mask := big.NewInt(0)

	wordsPhrase := strings.Fields(phrase)

	nbWords := len(wordsPhrase)

	// Nb words should be 12, 15, 18, 21 or 24
	if nbWords != 12 && nbWords != 15 && nbWords != 18 &&
		nbWords != 21 && nbWords != 24 {
		return "", "", newInvalidNumberOfPhraseWords()
	}

	// Read the wordlist file and extract words
	wlFile = filepath.Clean(wlFile)
	content, err := ioutil.ReadFile(wlFile)
	if err != nil {
		return "", "", newOpenWordlistError(wlFile)
	}
	// Split the file into words, handle multi words on one line
	wordsList := strings.Fields(string(content))

	// Load the word list in a map
	var wordsMap map[string]int = make(map[string]int)
	for idx, word := range wordsList {
		wordsMap[word] = idx
	}

	for _, wordP := range wordsPhrase {

		// Get the index of the word in the wordsMap/wordlist
		idx, found := wordsMap[wordP]
		if !found {
			return "", "", newWordNotFromTheWordlist(wordP, wlFile)
		}

		// Concatenate the index to find back the binary vector
		binary.BigEndian.PutUint16(wBytes[:], uint16(idx))

		bg = bg.Mul(bg, big.NewInt(2048)) // Shift 11 bits
		bg = bg.Or(bg, big.NewInt(0).SetBytes(wBytes[:]))
	}

	// The mask to get the checksum differs depending on the nb of words
	switch nbWords {
	case 12:
		mask = big.NewInt(15)
	case 15:
		mask = big.NewInt(31)
	case 18:
		mask = big.NewInt(63)
	case 21:
		mask = big.NewInt(127)
	case 24:
		mask = big.NewInt(255)
	}

	// Remove the checksum bits
	bg.Div(bg, big.NewInt(0).Add(mask, big.NewInt(1)))

	// The entropy is the rest of the bytes in the big int
	// The left bits are filled with 0
	entropy := bg.Bytes()
	offset := nbWords/12 - len(bg.Bytes())
	if offset > 0 {
		entropy := make([]byte, nbWords/12)
		copy(entropy[offset:], bg.Bytes())
	}

	// If the entropy is full 0, need to create the string manually
	// The []byte become [] by default
	var entropyStr string
	if string(entropy) == "" {
		entropyStr = string(strings.Repeat("0", 8*nbWords/3))
	} else {
		entropyStr = hex.EncodeToString(entropy)
	}

	// create the seed
	seedBytes, _ := phraseToSeed(phrase, passphrase)
	seed := hex.EncodeToString(seedBytes)

	return entropyStr, seed, nil
}

// VerifyPhraseAndSeed takes phrase, optional passphrase and a seed checks
// if the phrase generates the same seed using the passphrase.
func VerifyPhraseAndSeed(phrase, passphrase, seed string) (bool, error) {
	seedFromPhraseBytes, err := phraseToSeed(phrase, passphrase)
	seedFromPhrase := hex.EncodeToString(seedFromPhraseBytes)
	if err != nil {
		return false, err
	}
	if seedFromPhrase != seed {
		return false, nil
	}
	return true, nil
}
