package utils

import (
	"encoding/hex"
	"math/rand"
)

type ICodeGenerator interface {
	Generate(length int) (string, error)
}

type CodeGenerator struct{}

// GenerateUID generates a unique random hex value of length.
func (i CodeGenerator) Generate(length int) (uid string, err error) {
	// to reduce the amount of memory allocated divide the length by 2
	// as encoded hex string are of length * 2
	bytes := make([]byte, (length+1)/2)

	_, err = rand.Read(bytes)
	if err != nil {
		return "", err
	}

	uid = hex.EncodeToString(bytes)[:length]

	return uid, err
}
