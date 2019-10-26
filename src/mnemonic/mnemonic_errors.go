package mnemonic

import (
	"errors"
	"fmt"
)

// TODO give more descriptive error messages
const errEntropyNotDivisibleBy32 = "PA193 mnemonic: The 'entropy' bit-length '%d' is not divisible by 32."
const errEntropyIsNotHexadecimal = "PA193 mnemonic: The 'entropy' is not a hexadecimal string."
const errEntropyIsEmpty = "PA193 mnemonic: The 'entropy' is empty."
const errENTNotInRange = "PA193 mnemonic: The 'entropy' bit-length is not in the range %d-%d."
const errOpenWordlistError = "PA193_mnemonic: Could not open the wordlist file '%s'."

func newEntropyNotDivisibleBy32Error(length int) error {
	// FIXME add check for too big value
	// FIXME use special constructor for every error?
	return errors.New(fmt.Sprintf(errEntropyNotDivisibleBy32, length))
}

func newEntropyIsNotHexadecimalError() error {
	return errors.New(errEntropyIsNotHexadecimal)
}

func newEntropyIsEmptyError() error {
	return errors.New(errEntropyIsEmpty)
}

func newENTNotInRangeError() error {
	return errors.New(fmt.Sprintf(errENTNotInRange, lowerENTBound, upperENTBound))
}

func newOpenWordlistError(dictFilepath string) error {
	return errors.New(fmt.Sprintf(errOpenWordlistError, dictFilepath))
}
