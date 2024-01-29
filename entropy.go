package goregexpopulate

import "math/rand"

func entropyDefault(required bool) bool {
	e := entropyBase
	if !required {
		e += entropySkip
	}

	return rand.Intn(e) < entropyBase
}
