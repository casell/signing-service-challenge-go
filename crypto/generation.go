package crypto

import (
	"crypto"
	"strings"
)

var generators = make(map[string]Generator)

var normalizeAlgorithmName = strings.ToUpper

func registerGenerator(g Generator) {
	generators[normalizeAlgorithmName(g.Algorithm())] = g
}

// GetValidAlgorithms returns all known algorithms available for generation.
func GetValidAlgorithms() []string {
	list := make([]string, 0, len(generators))
	for k := range generators {
		list = append(list, k)
	}
	return list
}

// IsValidAlgorithm checks if a given algorithm is registered.
func IsValidAlgorithm(algorithmName string) bool {
	_, ok := generators[normalizeAlgorithmName(algorithmName)]
	return ok
}

// FromString returns a Generator using the given algorithm.
func FromString(algorithmName string) (Generator, error) {
	if gen, ok := generators[normalizeAlgorithmName(algorithmName)]; ok {
		return gen, nil
	} else {
		return nil, &ErrUnknownAlgorithm{algorithmName}
	}
}

// Generator is a common interface to create KeyPairs (either generating or loading)
type Generator interface {
	// Algorithm returns the algorithm name as a string.
	Algorithm() string
	// Generate generates a new KeyPair.
	Generate() (KeyPair, error)
	// Unmarshal loads a KeyPair from bytes.
	Unmarshal(priv []byte) (KeyPair, error)
}

// KeyPair is a common interface to RSA, ECC, ... keypairs
type KeyPair interface {
	// PrivateKey returns the keypair's private key as `crypto.Signer`.
	PrivateKey() crypto.Signer
	// PublicKey returns the keypair's public key.
	PublicKey() crypto.PublicKey
	// Equal verifies keypairs to be equals.
	Equal(KeyPair) bool
	// Marshal encodes the keypair to be written on disk.
	// It returns the public and the private key as a byte slice.
	Marshal() ([]byte, []byte, error)
}
