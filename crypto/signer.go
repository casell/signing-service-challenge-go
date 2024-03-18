package crypto

import (
	"crypto"
	"crypto/rand"
	"fmt"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

// GenericSigner represents a base Signer implementation.
type GenericSigner struct {
	cryptoSigner crypto.Signer
	cryptoHash   crypto.Hash
}

// NewGenericSigner returns a generic signer given a crypto.Signer and a crypto.Hash.
func NewGenericSigner(cryptoSigner crypto.Signer, cryptoHash crypto.Hash) (*GenericSigner, error) {
	if !cryptoHash.Available() {
		return nil, fmt.Errorf("hash function '%s' not available", cryptoHash.String())
	}
	return &GenericSigner{
		cryptoSigner: cryptoSigner,
		cryptoHash:   cryptoHash,
	}, nil
}

// Sign return the signature of the given data.
func (s *GenericSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hasher := s.cryptoHash.New()
	_, err := hasher.Write(dataToBeSigned)
	if err != nil {
		return nil, err
	}
	return s.cryptoSigner.Sign(rand.Reader, hasher.Sum(nil), s.cryptoHash)
}
