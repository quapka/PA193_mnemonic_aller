package mnemonic

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/bits"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func cleanInputEntropy(entropy string) ([]byte, error) {
	if len(entropy) == 0 {
		return nil, newEntropyIsEmptyError()
	}

	// FIXME test on huge inputs
	bytes, err := hex.DecodeString(entropy)
	if err != nil {
		return nil, newEntropyIsNotHexadecimalError()
	}

	//ENT := bits.Len(uint(bytes[0])) + (len(bytes)-1)*8
	ENT := getBinaryLength(bytes)
	inRange := lowerENTBound <= ENT && ENT <= upperENTBound
	if !inRange {
		return nil, newENTNotInRangeError()
	}

	//
	if ENT%32 != 0 {
		return nil, newEntropyNotDivisibleBy32Error(ENT)
	}
	return bytes, nil
}

func getBinaryLength(bytes []byte) int {
	return bits.Len(uint(bytes[0])) + (len(bytes)-1)*8
}

func calculateCheckSum(bytes []byte) string {
	hash := sha256.Sum256(bytes)
	ENT := getBinaryLength(bytes)
	checkSumLen := ENT / 32
	checkSum := hash[0] >> (8 - checkSumLen) // checkSumLen is always between 4-8
	// create formatting string like "%04s" - "%08s"
	checkSumFormat := fmt.Sprintf("%%0%ds", checkSumLen)
	checkSumBin := fmt.Sprintf(checkSumFormat, strconv.FormatInt(int64(checkSum), 2))
	return checkSumBin
}

func convertToBinary(bytes []byte) string {
	binary := ""
	for _, bin := range bytes {
		binary += fmt.Sprintf("%08s", strconv.FormatInt(int64(bin), 2))
	}
	return binary
}

func createGroups(binary string) (groups []string, err error) {
	length := len(binary)
	if length == 0 {
		return nil, nil
	}
	const groupSize = 11
	// FIXME use internal errors?
	// TODO maybe it does not need to be divisible?
	if length%groupSize != 0 {
		return nil, errors.New("'binary' length is not divisible by the group size")
	}
	// FIXME immediately create indices?
	for i := 0; i < (length / groupSize); i++ {
		group := binary[i*groupSize : (i+1)*groupSize]
		groups = append(groups, group)
	}
	// fmt.Println(groups)
	return groups, nil
}

func createIndices(groups []string) (indices []int64, err error) {
	for _, group := range groups {
		// FIXME handle errorneous case!
		ind, err := strconv.ParseInt(group, 2, 0)
		if err != nil {
			// FIXME better message and consistent
			return nil, errors.New("Cannot creat the phrase")
		}
		indices = append(indices, ind)

		// phrase += words[ind]
		// // FIXME remove trailing space
		// phrase += " "
	}
	return indices, nil
}

func createPhraseWords(indices []int64, words []string) (phrase []string, err error) {
	// FIXME perform input checking!
	for _, ind := range indices {
		phrase = append(phrase, words[ind])
	}
	return phrase, nil
}

func validateWordlist(wordList []string) (bool, error) {
	// assume it is a safe wordList
	valid := true
	// check duplicity
	frequency := make(map[string]int)

	for _, word := range wordList {
		_, exist := frequency[word]

		if exist {
			frequency[word] += 1
			// at least one duplicate found
			valid = false
		} else {
			frequency[word] = 1
		}
	}
	if !valid {
		// FIXME new error
		return false, errors.New("Contains duplicates!")
	}

	const expectedSize = 2048
	actualSize := len(wordList)
	if actualSize != expectedSize {
		// FIXME
		return false, errors.New("Not enough words in the wordlist")
	}
	return valid, nil
}

func loadWordlist(filepath string) ([]string, error) {
	dict, err := os.Open(filepath) // For read access.
	if err != nil {
		// TODO bubble the original error? Or simply in a wrapper?
		return nil, newOpenWordlistError(filepath)
	}
	defer dict.Close()
	// read the words
	scanner := bufio.NewScanner(dict)
	scanner.Split(bufio.ScanLines)
	var words []string
	for scanner.Scan() {
		word := cleanLine(scanner.Text())
		if !validateWord(word) {
			return nil, errors.New("Invalid word found")
		}
		// FIXME what does the line consist of?
		words = append(words, word)
	}
	// FIXME check for an error while reading
	// make sure the file is properly closed
	return words, nil
}

func cleanLine(line string) string {
	lower := strings.ToLower(line)
	trimmed := strings.TrimSpace(lower)
	return trimmed
}

func validateWord(word string) bool {
	wordPatten, err := regexp.Compile("^[a-z]+$")
	if err != nil {
		return false
	}
	return wordPatten.MatchString(word)
}