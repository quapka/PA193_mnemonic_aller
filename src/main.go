// Project: PA193_mnemonic_aller
// Maintainers UCO: 408788 497391 497577
// Description: Main file to test the API with command line

package main

import (
  "fmt"
  "mnemonic"
)


func main() {
  fmt.Println("Testing...")

  mnemonic.EntropyToPhraseAndSeed()


  // Entropy from phrase

  // phrase := "legal winner thank year wave sausage worth useful legal winner thank yellow"
  phrase := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
  wlfile := "wordlists/english.txt"

  if entropy, e := mnemonic.PhraseToEntropyAndSeed(phrase, wlfile); e != nil {
    fmt.Println(e)
  } else {
    fmt.Println("Entropy:", entropy)
  }




  mnemonic.VerifyPhraseAndSeed()

}

