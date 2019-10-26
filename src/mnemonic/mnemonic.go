// Project: PA193_mnemonic_aller
// Maintainers UCO: 408788, 497391, 497577
// Description: Mnemonic API

package mnemonic

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/bits"
	"os"
	"strconv"
)

const lowerENTBound = 128
const upperENTBound = 256

// entropy can be of various length, therefore it needs to a slice not an array
// FIXME check whether the underlying entropy array is changed, if so make a copy
// params:
// string entropy is a string of hexadecimal values
func EntropyToPhraseAndSeed(entropy string, dictFilepath string) (phrase, seed string, err error) {
	// FIXME test on huge inputs
	bytes, err := hex.DecodeString(entropy)
	if err != nil {
		return "", "", newEntropyIsNotHexadecimalError()
	}

	// FIXME refactor order of checks, check for emptiness first
	if len(bytes) == 0 {
		return "", "", newEntropyIsEmptyError()
	}

	ENT := bits.Len(uint(bytes[0])) + (len(bytes)-1)*8
	inRange := lowerENTBound <= ENT && ENT <= upperENTBound
	if !inRange {
		return "", "", newENTNotInRange()
	}

	//
	if ENT%32 != 0 {
		return "", "", newEntropyNotDivisibleBy32Error(ENT)
	}

	// checkSum
	// hash := sha256.New()
	// hash.Write(bytes)
	hash := sha256.Sum256(bytes)
	checkSumLen := ENT / 32
	checkSum := hash[0] >> (8 - checkSumLen) // checkSumLen is always between 4-8

	// create binary string
	binary := ""
	for _, bin := range bytes {
		binary += fmt.Sprintf("%08s", strconv.FormatInt(int64(bin), 2))
	}
	checkSumFormat := fmt.Sprintf("%%0%ds", checkSumLen)
	binary += fmt.Sprintf(checkSumFormat, strconv.FormatInt(int64(checkSum), 2))

	// create groups
	totalLen := ENT + checkSumLen
	const chunkSize = 11
	// FIXME  immediately crate indeces
	var chunks []string
	for i := 0; i < totalLen/chunkSize-1; i++ {
		chunk := binary[i*chunkSize : (i+1)*chunkSize]
		chunks = append(chunks, chunk)
	}

	dict, err := os.Open(dictFilepath) // For read access.
	if err != nil {
		// log.Fatal(err)
		// FIXME better message and consistent
		return "", "", errors.New("Cannot open the dictionary file")
	}
	// make sure the file is properly closed
	defer dict.Close()

	// FIXME check the dictionary file
	// create the phrase
	scanner := bufio.NewScanner(dict)
	scanner.Split(bufio.ScanLines)
	var words []string
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	phrase = ""
	for _, chunk := range chunks {
		// FIXME handle errorneous case!
		ind, err := strconv.ParseInt(chunk, 2, 0)
		if err != nil {
			// FIXME better message and consistent
			return "", "", errors.New("Cannot creat the phrase")
		}

		phrase += words[ind]
		// FIXME remove trailing space
		phrase += " "
	}
	fmt.Println(phrase)
	// FIXME add the implementation for the seed
	seed = ""

	return phrase, seed, nil
}

func PhraseToEntropyAndSeed() {
	// TODO
	fmt.Println("TODO entropy_to_mnemonic")
}

func VerifyPhraseAndSeed() {
	// TODO
	fmt.Println("TODO entropy_to_mnemonic")
}
