package main

import (
	"fmt"
	"math"
)

func GetBruteString(charset []rune, length int, counter uint32) (string, error) {
	n := uint32(len(charset))
	combinations := uint32(math.Pow(float64(n), float64(length)))
	if counter >= combinations {
		return "", fmt.Errorf("exceeded max counter limit for length [%d]: %d", length, (combinations - 1))
	}
	var pass []rune
	var quotient uint32 = counter
	for i := length - 1; i >= 0; i-- {
		pass = append([]rune{charset[quotient%n]}, pass...)
		quotient /= n
	}
	return string(pass), nil
}

func GetPasswordsForState(state State) ([]string, error) {
	var passwords []string
	var rangeError error = nil
	for i := state.IterationStart; i <= state.IterationEnd; i++ {
		password, err := GetBruteString(state.Charset, state.CurrentLength, i)
		if err != nil {
			rangeError = err
			break
		}
		passwords = append(passwords, password)
	}
	return passwords, rangeError
}
