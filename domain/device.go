package domain

import (
	"encoding/base64"
	"fmt"
	"sync"

	mycrypto "github.com/casell/signing-service-challenge/crypto"
	"github.com/google/uuid"
)

type Device struct {
	id                 uuid.UUID
	signatureAlgorithm string
	label              *string
	signatureCounter   uint
	lastSignatureB64   string
	signer             mycrypto.Signer
	keyPair            mycrypto.KeyPair
	lock               *sync.RWMutex
}

func (d *Device) ID() uuid.UUID {
	return d.id
}

func (d *Device) Label() *string {
	return d.label
}

func (d *Device) SignatureAlgorithm() string {
	return d.signatureAlgorithm
}

func (d *Device) KeyPair() mycrypto.KeyPair {
	return d.keyPair
}

func (d *Device) CounterAndLastSignature() (uint, string) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.signatureCounter, d.lastSignatureB64
}

func (d *Device) Sign(dataToBeSigned string) (string, string, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	extendedDataToBeSigned := fmt.Sprintf("%d_%s_%s", d.signatureCounter, dataToBeSigned, d.lastSignatureB64)
	signature, err := d.signer.Sign([]byte(extendedDataToBeSigned))
	if err != nil {
		return "", "", err
	}
	b64signature := base64.StdEncoding.EncodeToString(signature)

	d.signatureCounter++
	d.lastSignatureB64 = b64signature

	return b64signature, extendedDataToBeSigned, nil
}
