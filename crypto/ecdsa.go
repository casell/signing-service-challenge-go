package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

const (
	ECC_ALGORITHM_NAME = "ECC"
)

func init() {
	registerGenerator(&ECCGenerator{})
}

// ECCGenerator generates an ECC key pair.
type ECCGenerator struct{}

// Unmarshal loads an ECCKeyPair from bytes.
func (g *ECCGenerator) Unmarshal(priv []byte) (KeyPair, error) {
	return NewECCMarshaler().Decode(priv)
}

// Algorithm returns the algorithm name as a string.
func (g *ECCGenerator) Algorithm() string {
	return ECC_ALGORITHM_NAME
}

// Generate generates a new ECCKeyPair.
func (g *ECCGenerator) Generate() (KeyPair, error) {
	// Security has been ignored for the sake of simplicity.
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return &ECCKeyPair{
		Public:  &key.PublicKey,
		Private: key,
	}, nil
}

// ECCKeyPair is a DTO that holds ECC private and public keys.
type ECCKeyPair struct {
	Public  *ecdsa.PublicKey
	Private *ecdsa.PrivateKey
}

// PrivateKey returns the keypair's private key as `crytpo.Signer`.
func (k *ECCKeyPair) PrivateKey() crypto.Signer {
	return k.Private
}

// PublicKey returns the keypair's public key.
func (k *ECCKeyPair) PublicKey() crypto.PublicKey {
	return k.Public
}

// Equal verifies keypairs to be equals.
func (k *ECCKeyPair) Equal(x KeyPair) bool {
	return k.Private.Equal(x.PrivateKey())
}

// Marshal encodes the keypair to be written on disk.
// It returns the public and the private key as a byte slice.
func (k *ECCKeyPair) Marshal() ([]byte, []byte, error) {
	return NewECCMarshaler().Encode(*k)
}

// ECCMarshaler can encode and decode an ECC key pair.
type ECCMarshaler struct{}

// NewECCMarshaler creates a new ECCMarshaler.
func NewECCMarshaler() ECCMarshaler {
	return ECCMarshaler{}
}

// Encode takes an ECCKeyPair and encodes it to be written on disk.
// It returns the public and the private key as a byte slice.
func (m ECCMarshaler) Encode(keyPair ECCKeyPair) ([]byte, []byte, error) {
	privateKeyBytes, err := x509.MarshalECPrivateKey(keyPair.Private)
	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(keyPair.Public)
	if err != nil {
		return nil, nil, err
	}

	encodedPrivate := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE_KEY",
		Bytes: privateKeyBytes,
	})

	encodedPublic := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC_KEY",
		Bytes: publicKeyBytes,
	})

	return encodedPublic, encodedPrivate, nil
}

// Decode assembles an ECCKeyPair from an encoded private key.
func (m ECCMarshaler) Decode(privateKeyBytes []byte) (*ECCKeyPair, error) {
	//XXX: nil block segfaulted
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, errors.New("no PEM data found")
	}
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &ECCKeyPair{
		Private: privateKey,
		Public:  &privateKey.PublicKey,
	}, nil
}
