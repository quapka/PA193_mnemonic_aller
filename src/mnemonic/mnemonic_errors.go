package mnemonic

import (
	"errors"
	"fmt"
)

// TODO give more descriptive error messages
const errEntropyNotDivisibleBy32 = "The 'entropy' bit-length '%d' is not divisible by 32."
const errEntropyIsNotHexadecimal = "The 'entropy' is not a hexadecimal string."
const errEntropyIsEmpty = "The 'entropy' is empty."
const errENTNotInRange = "The 'entropy' bit-length is not in the range %d-%d."

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

func newENTNotInRange() error {
	return errors.New(fmt.Sprintf(errENTNotInRange, lowerENTBound, upperENTBound))
}
