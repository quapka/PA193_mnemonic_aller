// Project: PA193_mnemonic_aller
// Maintainers UCO: 408788, 497391, 497577
// Description: Mnemonic API

package mnemonic

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)

// FIXME add documention of functions
// entropy can be of various length, therefore it needs to a slice not an array
// FIXME check whether the underlying entropy array is changed, if so make a copy
// params:
// string entropy is a string of hexadecimal values
func EntropyToPhraseAndSeed(entropy, passphrase, dictFilepath string) (phrase, seed string, err error) {
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
	// create the mnemonic phrase
	phraseWords, _ := createPhraseWords(indices, wordList)
	phrase = strings.Join(phraseWords, " ")
	// create the seed value
	seedBytes, _ := phraseToSeed(phrase, passphrase)
	seed = hex.EncodeToString(seedBytes)

	return phrase, seed, nil
}

// FIXME make wlfile naming consistetn!
func PhraseToEntropyAndSeed(phrase, passphrase, wlfile string) (string, string, error) {

	var wBytes [2]byte

	// Use a big int for arbitrary-precision arithmetic of big numbers
	bg := big.NewInt(0)
	checksum := big.NewInt(0)
	mask := big.NewInt(0)

	wordsPhrase := strings.Fields(phrase)

	nbWords := len(wordsPhrase)

	// Nb words should be 12, 15, 18, 21 or 24
	if nbWords != 12 && nbWords != 15 && nbWords != 18 && nbWords != 21 && nbWords != 24 {
		return "", "", errors.New("Phrase Invalid")
	}

	// Read the wordlist file and extract words
	content, e := ioutil.ReadFile(wlfile)
	if e != nil {
		return "", "", e
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
		if found == false {
			return "", "", errors.New("Phrase word not in wordlist: " + wordP)
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

	// Get the checksum
	checksum = checksum.And(bg, mask)
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

func VerifyPhraseAndSeed(phrase_to_verify, passphrase, seed_to_verify, wlfile string) int {
	var seed string
	// FIXME
	_, seed, _ = PhraseToEntropyAndSeed(phrase_to_verify, passphrase, wlfile) /* WARNING NEED THE SEED IN THE RETURN NOT ONLY THE ERROR */
	if seed_to_verify == seed {
		fmt.Println("The phrase and the seed correspond !!")
		return 0
	} else {
		fmt.Println("The phrase and the seed do NOT correspond !!")
		return -1
	}
	return -1
}

// FIXME lowercase
func PBKDF2_SHA512_F(password, salt []byte, count, l_counter int) ([]byte, int) {

	U_1_to_c := make([][]byte, count)
	if U_1_to_c == nil {
		return nil, -1
	}

	var INT_i_l [4]byte
	INT_i_l[0] = byte((l_counter >> 24))
	INT_i_l[1] = byte((l_counter >> 16)) /* INT (i) is a four-octet encoding of the integer i, most significant octet first. */
	INT_i_l[2] = byte((l_counter >> 8))
	INT_i_l[3] = byte((l_counter))

	U_1_to_c[0] = make([]byte, 0, 64) /* U_1 = PRF (P, S || INT (i)) , */
	if U_1_to_c[0] == nil {
		return nil, -1
	}

	sha_512 := hmac.New(sha512.New, password)
	if sha_512 == nil {
		return nil, -1
	}

	sha_512.Write(salt)
	sha_512.Write(INT_i_l[:4])

	U_1_to_c[0] = sha_512.Sum(U_1_to_c[0])
	if U_1_to_c[0] == nil {
		return nil, -1
	}

	sha_512.Reset()

	for i := 1; i < count; i++ { /* U_2 = PRF (P, U_1) , */
		sha_512.Reset() /* ... */

		U_1_to_c[i] = make([]byte, 64) /* U_c = PRF (P, U_{c-1}) . */
		if U_1_to_c[i] == nil {
			return nil, -1
		}

		sha_512.Write(U_1_to_c[i-1])

		U_1_to_c[i] = sha_512.Sum(nil)
		if U_1_to_c[i] == nil {
			return nil, -1
		}

		sha_512.Reset()
	}

	output := make([]byte, 64)
	if output == nil {
		return nil, -1
	}

	output = U_1_to_c[0]

	for i := 1; i < count; i++ { /* F (P, S, c, i) = U_1 \xor U_2 \xor ... \xor U_c */
		for j := range U_1_to_c[i] {
			output[j] ^= U_1_to_c[i][j]
		}
	}
	return output, 0
}

/* https://www.ietf.org/rfc/rfc2898.txt

Translation variable with RFC :
U_1, U_2, U_3 ... U_c is the array U_1_to_c, begin at 0 end at c-1
T_1, T_2, T_3 ... T_l is the array T_1_to_l, begin at 0 end at l-1
hLen is hLen
dkLen is output_len
P is password
S is salt
c is count
INT(i) is INT_i_l
l_counter in the program is the index until l in the RFC
*/
// FIXME move to utils.go
func PBKDF2_SHA512(password, salt []byte, count, output_len int) ([]byte, int) {
	if output_len != 64 { /* Length of SHA-512 */ /* 1. If dkLen > (2^32 - 1) * hLen, output "derived key too long" and stop.*/
		return nil, -1
	} else {

		err := 0
		hLen := 64 /* Length of SHA-512 */
		var l int

		if hLen != 0 {
			l = output_len / hLen /* Should be equal to 1 !*/ /* l = CEIL (dkLen / hLen) , */
		} else {
			return nil, -1
		}
		// r := output_len -(l-1)*hLen        /* Should be equal to output_len, so 64 bytes */  /* r = dkLen - (l - 1) * hLen . */
		/* Commented because it is an unused variable */

		T_1_to_l := make([][]byte, output_len)
		if T_1_to_l == nil {
			return nil, -1
		}

		for i := 0; i < l; i++ { /* T_1 = F (P, S, c, 1) ,*/
			T_1_to_l[i] = make([]byte, output_len) /* T_2 = F (P, S, c, 2) ,*/
			if T_1_to_l[i] == nil {
				return nil, -1
			}

			T_1_to_l[i], err = PBKDF2_SHA512_F(password, salt, count, i+1) /* i+1 because begin l in RFC        ...         */
			if err < 0 {
				return nil, -1
			}
		} /* T_l = F (P, S, c, l) , */

		output := make([]byte, output_len)
		if output == nil {
			return nil, -1
		}

		output = T_1_to_l[0]

		/* This part is only used if output_len is  greater than SHA512. (64 bytes)
		* In bip39, the output_len is always 64 bytes, that is why we comment this part of code
		 */
		// for i := 1; i < l; i++ {
		// 	output = append(output, T_1_to_l[i]...) /* DK = T_1 || T_2 ||  ...  || T_l<0..r-1> */
		// }
		return output, 0
	}
}

/* This function converts a mnemonic phrase to the corresponding seed using PBKDF2. */
// FIXME move to utils
func phraseToSeed(phrase, passphrase string) (seed []byte, err int) {
	seed, err = PBKDF2_SHA512([]byte(phrase), []byte("mnemonic"+passphrase), 2048, 64)
	if err < 0 {
		fmt.Fprintf(os.Stderr, "Error in PBKDF2_SHA512")
		seed = nil
	}
	return seed, err
}
