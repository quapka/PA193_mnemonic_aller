// Project: PA193_mnemonic_aller
// Mainteners UCO: 408788 497391 497577
// Description: Main file to test the API with command line

package main

import (
  "fmt"
  "./mnemonic"
  "encoding/hex"
)


func main() {
  fmt.Println("Testing...")

  mnemonic.EntropyToPhraseAndSeed()
  mnemonic.PhraseToEntropyAndSeed("phrase")
  mnemonic.VerifyPhraseAndSeed("phrase","phrase")

  tmp,_ := mnemonic.PhraseToSeed("void come effort suffer camp survey warrior heavy shoot primary clutch crush open amazing screen patrol group space point ten exist slush involve unfold","TREZOR")
  fmt.Println(hex.EncodeToString(tmp))
}
