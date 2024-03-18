package domain

import (
	mycrypto "github.com/casell/signing-service-challenge/crypto"
	"github.com/google/uuid"
)

type SigningDevice interface {
	ID() uuid.UUID
	SignatureAlgorithm() string
	Label() *string
	KeyPair() mycrypto.KeyPair
	CounterAndLastSignature() (uint, string)
	Sign(dataToBeSigned string) (string, string, error)
}

type SigningDeviceFactory interface {
	New(signatureAlgorithm string, label *string) (SigningDevice, error)
}
