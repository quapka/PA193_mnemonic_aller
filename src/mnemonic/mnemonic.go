// Project: PA193_mnemonic_aller
// Maintainers UCO: 408788, 497391, 497577
// Description: Mnemonic API

package mnemonic

import (
  "fmt"
  "strings"
  "errors"
  "io/ioutil"
  "encoding/binary"
  "math/big"
  "encoding/hex"
)

func EntropyToPhraseAndSeed() {
  // TODO
  fmt.Println("TODO entropy_to_mnemonic")
}





func PhraseToEntropyAndSeed(phrase string, wlfile string) (string, error) {

  var wBytes [2]byte

  // Use a big int for arbitrary-precision arithmetic of big numbers
  bg := big.NewInt(0)
  checksum := big.NewInt(0)
  mask := big.NewInt(0)

  wordsPhrase := strings.Fields(phrase)

  nbWords := len(wordsPhrase)

  // Nb words should be 12, 15, 18, 21 or 24
  if nbWords!=12 && nbWords!=15 && nbWords!=18 && nbWords!=21 && nbWords!=24 {
    return "", errors.New("Phrase Invalid")
  }

  // Read the wordlist file and extract words
  content, e := ioutil.ReadFile(wlfile)
  if e != nil {
    return "", e
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
      return "", errors.New("Phrase word not in wordlist: " + wordP)
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

  return entropyStr, nil
}



func VerifyPhraseAndSeed() {
  // TODO
  fmt.Println("TODO entropy_to_mnemonic")
}

