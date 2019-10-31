package mnemonic

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const lowerENTBound = 128
const upperENTBound = 256

func cleanInputEntropy(entropy string) ([]byte, error) {
	if len(entropy) == 0 {
		return nil, newEntropyIsEmptyError()
	}

	// FIXME test on huge inputs
	bytes, err := hex.DecodeString(entropy)
	if err != nil {
		return nil, newEntropyIsNotHexadecimalError()
	}

	ENT := getBinaryLength(bytes)
	notInRange := !(lowerENTBound <= ENT && ENT <= upperENTBound)
	if notInRange {
		return nil, newENTNotInRangeError()
	}

	if ENT%32 != 0 {
		return nil, newEntropyNotDivisibleBy32Error(ENT)
	}
	return bytes, nil
}

func getBinaryLength(bytes []byte) int {
	return len(bytes) * 8
}

func calculateCheckSum(bytes []byte) string {
	hash := sha256.Sum256(bytes)
	ENT := getBinaryLength(bytes)
	checkSumLen := ENT / 32
	// checkSumLen is always between 4-8
	checkSum := hash[0] >> (8 - checkSumLen)
	// create formatting string like "%04s" - "%08s"
	checkSumFormat := fmt.Sprintf("%%0%ds", checkSumLen)
	checkSumBin := fmt.Sprintf(checkSumFormat,
		strconv.FormatInt(int64(checkSum), 2))
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
		return nil, newBinaryLenghtIsNotDivisibleByGroupSize()
	}
	// FIXME immediately create indices?
	for i := 0; i < (length / groupSize); i++ {
		group := binary[i*groupSize : (i+1)*groupSize]
		groups = append(groups, group)
	}
	return groups, nil
}

func createIndices(groups []string) (indices []int64, err error) {
	for _, group := range groups {
		// FIXME handle errorneous case!
		ind, err := strconv.ParseInt(group, 2, 0)
		if err != nil {
			// FIXME better message and consistent
			return nil, newCannotParseIntegerError(group)
		}
		indices = append(indices, ind)
	}
	return indices, nil
}

func createPhraseWords(indices []int64,
	words []string) (phraseWords []string, err error) {
	// FIXME perform input checking!
	for _, ind := range indices {
		phraseWords = append(phraseWords, words[ind])
	}
	return phraseWords, nil
}

func validateWordlist(wordList []string) error {
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
		return newWordlistContainsDuplicatesError()
	}

	const expectedSize = 2048
	actualSize := len(wordList)
	if actualSize != expectedSize {
		return newNotExpectedWordlistSizeError()
	}
	return nil
}

func loadWordlist(wlFilePath string) ([]string, error) {
	wlFilePath = filepath.Clean(wlFilePath)
	wlFile, err := os.Open(wlFilePath) // For read access.
	if err != nil {
		// TODO bubble the original error? Or simply in a wrapper?
		return nil, newOpenWordlistError(wlFilePath)
	}
	defer wlFile.Close()

	scanner := bufio.NewScanner(wlFile)
	scanner.Split(bufio.ScanLines)
	var words []string

	for scanner.Scan() {
		word := cleanLine(scanner.Text())
		if !validateWord(word) {
			return nil, newInvalidWordError(word)
		}
		words = append(words, word)
	}
	if err := validateWordlist(words); err != nil {
		return nil, err
	}
	return words, nil
}

func cleanLine(line string) string {
	lower := strings.ToLower(line)
	trimmed := strings.TrimSpace(lower)
	return trimmed
}

func validateWord(word string) bool {
	wordPatten, err := regexp.Compile(`^[^\s]+$`)
	if err != nil {
		return false
	}
	return wordPatten.MatchString(word)
}

// Pbkdf2Sha512F
// This function is the implementation of the function F in Pbkdf2
// according to the RFC 2898 notation
// https://www.ietf.org/rfc/rfc2898.txt
// This function is called by the Pbkdf2Sha512 function and should not
// be used in another context
func Pbkdf2Sha512F(password, salt []byte, count, lCounter int) ([]byte, int) {
	// Translation variable with RFC :
	// U_1, U_2, U_3 ... U_c is the array U1ToC, begin at 0 end at c-1
	// T_1, T_2, T_3 ... T_l is the array T1Tol, begin at 0 end at l-1
	// hLen is hLen
	// dkLen is OutputLen
	// P is password
	// S is salt
	// c is count
	// INT(i) is INTil
	// lCounter in the program is the index until l in the RFC

	U1ToC := make([][]byte, count)

	var INTil [4]byte
	INTil[0] = byte((lCounter >> 24))
	INTil[1] = byte((lCounter >> 16)) /* INT (i) is a four-octet encoding of the integer i, most significant octet first. */
	INTil[2] = byte((lCounter >> 8))
	INTil[3] = byte((lCounter))

	U1ToC[0] = make([]byte, 0, 64) /* U_1 = PRF (P, S || INT (i)) , */

	Sha512 := hmac.New(sha512.New, password)
	if Sha512 == nil {
		return nil, -1
	}

	Sha512.Write(salt)
	Sha512.Write(INTil[:4])

	U1ToC[0] = Sha512.Sum(U1ToC[0])
	if U1ToC[0] == nil {
		return nil, -1
	}

	Sha512.Reset()

	for i := 1; i < count; i++ { /* U_2 = PRF (P, U_1) , */
		Sha512.Reset() /* ... */

		U1ToC[i] = make([]byte, 64) /* U_c = PRF (P, U_{c-1}) . */

		Sha512.Write(U1ToC[i-1])

		U1ToC[i] = Sha512.Sum(nil)
		if U1ToC[i] == nil {
			return nil, -1
		}

		Sha512.Reset()
	}

	output := make([]byte, 64)

	output = U1ToC[0]

	for i := 1; i < count; i++ { /* F (P, S, c, i) = U_1 \xor U_2 \xor ... \xor U_c */
		for j := range U1ToC[i] {
			output[j] ^= U1ToC[i][j]
		}
	}
	return output, 0
}

// Pbkdf2Sha512
// Implementation of Pbkdf2Sha512 according to the RFC 2898
// https://www.ietf.org/rfc/rfc2898.txt
// Parameter :
// password   : is the password that will be derived (P in RFC)
// salt 		  : is the salt that will be added to password (S in RFC)
// count 		  : Number of iteration of SHA-512 (c in the RFC)
// OutputLen : Length of the derived password (output) MUST BE 64 as 64 bytes
func Pbkdf2Sha512(password, salt []byte, count, OutputLen int) ([]byte, int) {
	// Translation variable with RFC :
	// U_1, U_2, U_3 ... U_c is the array U1ToC, begin at 0 end at c-1
	// T_1, T_2, T_3 ... T_l is the array T1Tol, begin at 0 end at l-1
	// hLen is hLen
	// dkLen is OutputLen
	// P is password
	// S is salt
	// c is count
	// INT(i) is INTil
	// lCounter in the program is the index until l in the RFC
	if OutputLen != 64 { /* Length of SHA-512 */ /* 1. If dkLen > (2^32 - 1) * hLen, output "derived key too long" and stop.*/
		return nil, -1
	} else {
		err := 0
		hLen := 64 /* Length of SHA-512 */
		var l int

		if hLen != 0 {
			l = OutputLen / hLen /* Should be equal to 1 !*/ /* l = CEIL (dkLen / hLen) , */
		} else {
			return nil, -1
		}
		// r := OutputLen -(l-1)*hLen        /* Should be equal to OutputLen, so 64 bytes */  /* r = dkLen - (l - 1) * hLen . */
		/* Commented because it is an unused variable */

		T1Tol := make([][]byte, OutputLen)

		for i := 0; i < l; i++ { /* T_1 = F (P, S, c, 1) ,*/
			T1Tol[i] = make([]byte, OutputLen) /* T_2 = F (P, S, c, 2) ,*/

			T1Tol[i], err = Pbkdf2Sha512F(password, salt, count, i+1) /* i+1 because begin l in RFC        ...         */
			if err < 0 {
				return nil, -1
			}
		} /* T_l = F (P, S, c, l) , */

		output := T1Tol[0]

		/* This part is only used if OutputLen is  greater than SHA512. (64 bytes)
		* In bip39, the OutputLen is always 64 bytes, that is why we comment this part of code
		 */
		// for i := 1; i < l; i++ {
		// 	output = append(output, T1Tol[i]...) /* DK = T_1 || T_2 ||  ...  || T_l<0..r-1> */
		// }
		return output, 0
	}
}
