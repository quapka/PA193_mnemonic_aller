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

// TODO remove _
const errOpenWordlistError = "PA193 mnemonic: Could not open the wordlist file '%s'."
const errInvalidWordError = "PA193 mnemonic: Wordlist contains an invalid word '%s'."
const errCannotParseIntegerError = "PA193 mnemonic: Cannot parse '%s' as an integer."
const errWordlistContainsduplicates = "PA193 mnemonic: The wordlist provided contains duplicate words. Make sure it is unique."
const errNotExpectedWordlistSize = "PA193 mnemonic: The wordlist size is not 2048."

func newEntropyNotDivisibleBy32Error(length int) error {
	// FIXME add check for too big value
	// FIXME use special constructor for every error?
	return fmt.Errorf(errEntropyNotDivisibleBy32, length)
}

func newEntropyIsNotHexadecimalError() error {
	return errors.New(errEntropyIsNotHexadecimal)
}

func newEntropyIsEmptyError() error {
	return errors.New(errEntropyIsEmpty)
}

func newENTNotInRangeError() error {
	return fmt.Errorf(errENTNotInRange, lowerENTBound, upperENTBound)
}

func newOpenWordlistError(dictFilepath string) error {
	return fmt.Errorf(errOpenWordlistError, dictFilepath)
}

func newInvalidWordError(word string) error {
	return fmt.Errorf(errInvalidWordError, word)
}

func newCannotParseIntegerError(integer string) error {
	return fmt.Errorf(errCannotParseIntegerError, integer)
}

func newWordlistContainsDuplicatesError() error {
	return errors.New(errWordlistContainsduplicates)
}

func newNotExpectedWordlistSizeError() error {
	return errors.New(errNotExpectedWordlistSize)
}
