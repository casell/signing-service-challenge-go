package domain

import (
	"crypto"
	"encoding/base64"
	"sync"

	mycrypto "github.com/casell/signing-service-challenge/crypto"
	"github.com/google/uuid"
)

type DefaultDeviceFactory struct{}

func NewDefaultDeviceFactory() *DefaultDeviceFactory {
	return &DefaultDeviceFactory{}
}

func (f *DefaultDeviceFactory) Restore(id uuid.UUID, signatureAlgorithm string, label *string, signatureCounter uint, lastSignatureB64 string, keyPair mycrypto.KeyPair, lock *sync.RWMutex) (SigningDevice, error) {
	if ok := mycrypto.IsValidAlgorithm(signatureAlgorithm); !ok {
		return nil, &ErrInvalidAlgorithm{signatureAlgorithm}
	}
	s, err := mycrypto.NewGenericSigner(keyPair.PrivateKey(), crypto.SHA256)
	if err != nil {
		return nil, err
	}
	return &Device{
		id:                 id,
		signatureAlgorithm: signatureAlgorithm,
		label:              label,
		signatureCounter:   signatureCounter,
		lastSignatureB64:   lastSignatureB64,
		signer:             s,
		keyPair:            keyPair,
		lock:               lock,
	}, nil
}

func (*DefaultDeviceFactory) New(signatureAlgorithm string, label *string) (SigningDevice, error) {
	g, err := mycrypto.FromString(signatureAlgorithm)
	if err != nil {
		return nil, &ErrInvalidAlgorithm{signatureAlgorithm}
	}

	kp, err := g.Generate()
	if err != nil {
		return nil, err
	}

	uniqueId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	s, err := mycrypto.NewGenericSigner(kp.PrivateKey(), crypto.SHA256)
	if err != nil {
		return nil, err
	}

	d := &Device{
		id:                 uniqueId,
		signatureAlgorithm: g.Algorithm(),
		label:              label,
		keyPair:            kp,
		signer:             s,
		signatureCounter:   0,
		lastSignatureB64:   base64.StdEncoding.EncodeToString([]byte(uniqueId.String())),
		lock:               &sync.RWMutex{},
	}
	return d, nil
}
