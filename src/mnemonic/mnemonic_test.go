// Project: PA193_mnemonic_aller
// Mainteners UCO: 408788 497391 497577
// Description: Mnemonic API unit testing

package mnemonic

import (
  "testing"
)


func TestEntropyToPhraseAndSeed(t *testing.T) {
  a := 4
  if 1 != a {
    t.Errorf("Here the reason why and variables %q", a)
  }
}
