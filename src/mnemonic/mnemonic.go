// Project: PA193_mnemonic_aller
// Maintainers UCO: 408788, 497391, 497577
// Description: Mnemonic API

package mnemonic

import (
  "fmt"
)

func EntropyToPhraseAndSeed() {
  // TODO
  fmt.Println("TODO entropy_to_mnemonic")
}

func PhraseToEntropyAndSeed(phrase string) (entropy string,seed string){
  // TODO
  fmt.Println("TODO entropy_to_mnemonic")
  // Trasnformation of phrase to entropy and seed
  seed = phrase     // For test, must be deleted and replaced
  return entropy, seed
}

func VerifyPhraseAndSeed(phrase_to_verify,seed_to_verify string) int{
  // TODO
  // fmt.Println("TODO entropy_to_mnemonic")
  var seed string
  _, seed = PhraseToEntropyAndSeed(phrase_to_verify)
  if seed_to_verify==seed{
    fmt.Println("The phrase and the seed correspond !!")
    return 0
  } else {
    fmt.Println("The phrase and the seed do NOT correspond !!")
    return -1
  }
  return -1
}
