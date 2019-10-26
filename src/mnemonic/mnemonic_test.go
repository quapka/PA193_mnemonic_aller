// Project: PA193_mnemonic_aller
// Mainteners UCO: 408788 497391 497577
// Description: Mnemonic API unit testing

package mnemonic

import (
	//"errors"
	"fmt"
	"testing"
)

func gotExp(got, exp string) string {
	return fmt.Sprintf("\nGot: %s\nExp: %s\n", got, exp)
}

func TestEntropyIsNotEmpty(t *testing.T) {
	entropy := ""
	phrase, seed, err := EntropyToPhraseAndSeed(entropy, "english.txt")

	expectedPhrase := ""
	expectedSeed := ""
	expectedErr := newEntropyIsEmptyError()

	if err.Error() != expectedErr.Error() {
		t.Error(gotExp(err.Error(), expectedErr.Error()))
	}
	if phrase != expectedPhrase {
		t.Error(gotExp(phrase, expectedPhrase))
	}
	if seed != expectedSeed {
		t.Error(gotExp(seed, expectedSeed))
	}
}

func TestENTIsInRange(t *testing.T) {
	// FIXME smaller, in range, higher, catch off by one errors
	// fairly small input
	entropy := "FF"
	expectedErr := newENTNotInRangeError()
	// FIXME
	_, _, err := EntropyToPhraseAndSeed(entropy, "english.txt")
	if err.Error() != expectedErr.Error() {
		t.Error(gotExp(err.Error(), expectedErr.Error()))
	}
	// smaller by 1
	entropy = "7FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF" // 2**127 - 1
	_, _, err = EntropyToPhraseAndSeed(entropy, "english.txt")
	if err.Error() != expectedErr.Error() {
		t.Error(gotExp(err.Error(), expectedErr.Error()))
	}
	// bigger by 1
	// FIXME different output this and Go Playground version, difference between 01
	// or 1 in the beginning
	entropy = "010000000000000000000000000000000000000000000000000000000000000000" // 2 ** 256
	_, _, err = EntropyToPhraseAndSeed(entropy, "english.txt")
	if err.Error() != expectedErr.Error() {
		t.Error(gotExp(err.Error(), expectedErr.Error()))
	}
}

func TestENTIsMultipleOf32x(t *testing.T) {
	entropy := "0800000000000000000000000000000000000000000000"
	_, _, err := EntropyToPhraseAndSeed(entropy, "english.txt")
	expectedErr := newEntropyNotDivisibleBy32Error(180)
	if err.Error() != expectedErr.Error() {
		t.Error(gotExp(err.Error(), expectedErr.Error()))
	}
}

func TestEntropyIsHexadecimal(t *testing.T) {
	entropy := "XX"
	phrase, seed, err := EntropyToPhraseAndSeed(entropy, "english.txt")
	expectedPhrase := ""
	expectedSeed := ""
	expectedErr := newEntropyIsNotHexadecimalError()
	// Comparing the actual error string
	if err.Error() != expectedErr.Error() {
		// if errors.Is(err, expectedErr) {
		t.Error(fmt.Sprintf("\nGot: %s\nExp: %s\n", err, expectedErr))
	}
	if phrase != expectedPhrase {
		t.Error(gotExp(phrase, expectedPhrase))
	}

	if seed != expectedSeed {
		t.Error(gotExp(seed, expectedSeed))
	}
}

func TestCannotOpenWordlistFile(t *testing.T) {
	entropy := "B7CB8EE904628CEC2B6779C0FB8B1B91" // 2**127 - 1
	phrase, seed, err := EntropyToPhraseAndSeed(entropy, "does/not/exist")
	expectedPhrase := ""
	expectedSeed := ""
	expectedErr := newOpenWordlistError("does/not/exist")

	if err.Error() != expectedErr.Error() {
		t.Error(gotExp(string(err.Error()), string(expectedErr.Error())))
	}
	if phrase != expectedPhrase {
		t.Error(gotExp(phrase, expectedPhrase))
	}

	if seed != expectedSeed {
		t.Error(gotExp(seed, expectedSeed))
	}
}
