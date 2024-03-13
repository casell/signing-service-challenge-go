package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

const (
	RSA_ALGORITHM_NAME = "RSA"
)

func init() {
	registerGenerator(&RSAGenerator{})
}

// RSAGenerator generates a RSA key pair.
type RSAGenerator struct{}

// Algorithm returns the algorithm name as a string.
func (g *RSAGenerator) Algorithm() string {
	return RSA_ALGORITHM_NAME
}

// Generate generates a new RSAKeyPair.
func (g *RSAGenerator) Generate() (KeyPair, error) {
	// Security has been ignored for the sake of simplicity.
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		return nil, err
	}

	return &RSAKeyPair{
		Public:  &key.PublicKey,
		Private: key,
	}, nil
}

// Unmarshal loads an RSAKeyPair from bytes.
func (g *RSAGenerator) Unmarshal(priv []byte) (KeyPair, error) {
	return NewRSAMarshaler().Unmarshal(priv)
}

// RSAKeyPair is a DTO that holds RSA private and public keys.
type RSAKeyPair struct {
	Public  *rsa.PublicKey
	Private *rsa.PrivateKey
}

// PrivateKey returns the keypair's private key as `crypto.Signer`.
func (k *RSAKeyPair) PrivateKey() crypto.Signer {
	return k.Private
}

// PublicKey returns the keypair's public key.
func (k *RSAKeyPair) PublicKey() crypto.PublicKey {
	return k.Public
}

// Equal verifies keypairs to be equals.
func (k *RSAKeyPair) Equal(x KeyPair) bool {
	return k.Private.Equal(x.PrivateKey())
}

// Marshal encodes the keypair to be written on disk.
// It returns the public and the private key as a byte slice.
func (k *RSAKeyPair) Marshal() ([]byte, []byte, error) {
	return NewRSAMarshaler().Marshal(*k)
}

// RSAMarshaler can encode and decode an RSA key pair.
type RSAMarshaler struct{}

// NewRSAMarshaler creates a new RSAMarshaler.
func NewRSAMarshaler() RSAMarshaler {
	return RSAMarshaler{}
}

// Marshal takes an RSAKeyPair and encodes it to be written on disk.
// It returns the public and the private key as a byte slice.
func (m RSAMarshaler) Marshal(keyPair RSAKeyPair) ([]byte, []byte, error) {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(keyPair.Private)
	publicKeyBytes := x509.MarshalPKCS1PublicKey(keyPair.Public)

	encodedPrivate := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA_PRIVATE_KEY",
		Bytes: privateKeyBytes,
	})

	encodePublic := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA_PUBLIC_KEY",
		Bytes: publicKeyBytes,
	})

	return encodePublic, encodedPrivate, nil
}

// Unmarshal takes an encoded RSA private key and transforms it into a rsa.PrivateKey.
func (m RSAMarshaler) Unmarshal(privateKeyBytes []byte) (*RSAKeyPair, error) {
	//XXX: nil block segfaulted
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, errors.New("no PEM data found")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &RSAKeyPair{
		Private: privateKey,
		Public:  &privateKey.PublicKey,
	}, nil
}
