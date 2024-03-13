package api

import (
	"context"
	"errors"
	"testing"

	"github.com/casell/signing-service-challenge/domain"
	"github.com/casell/signing-service-challenge/generated/signingapi"
	mockCrypto "github.com/casell/signing-service-challenge/mocks/github.com/casell/signing-service-challenge/crypto"
	mockDomain "github.com/casell/signing-service-challenge/mocks/github.com/casell/signing-service-challenge/domain"
	mockPersistence "github.com/casell/signing-service-challenge/mocks/github.com/casell/signing-service-challenge/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateDeviceUnableToAdd(t *testing.T) {
	algo := string(signingapi.DeviceRequestSignatureAlgorithmRSA)

	var label *string = nil

	addErr := errors.New("Add error")

	mockDevice := mockDomain.NewMockSigningDevice(t)
	mockFactory := setupMockFactory(t, algo, label, mockDevice)
	mockStorage := mockPersistence.NewMockStorage(t)
	mockStorage.EXPECT().Add(mockDevice).Return(addErr)

	dh := NewDeviceHandler(mockStorage, mockFactory)

	reqLabel := signingapi.OptNilString{}

	res, err := dh.CreateDevice(context.TODO(), &signingapi.DeviceRequest{
		SignatureAlgorithm: signingapi.DeviceRequestSignatureAlgorithmRSA,
		Label:              reqLabel,
	})

	assert.Nil(t, res)
	if assert.Error(t, err) {
		assert.Equal(t, addErr, err)
	}
}

func setupMockDevice(t *testing.T, id uuid.UUID, algo, pub, priv, signature string, counter int, label *string) domain.SigningDevice {
	mockDevice := mockDomain.NewMockSigningDevice(t)
	mockDevice.EXPECT().ID().Return(id)
	mockKP := mockCrypto.NewMockKeyPair(t)
	mockKP.EXPECT().Marshal().Return([]byte(pub), []byte(priv), nil)
	mockDevice.EXPECT().KeyPair().Return(mockKP)
	mockDevice.EXPECT().CounterAndLastSignature().Return(uint(counter), signature)
	mockDevice.EXPECT().SignatureAlgorithm().Return(algo)
	mockDevice.EXPECT().Label().Return(label)
	return mockDevice
}

func setupMockFactory(t *testing.T, algo string, label *string, mockDevice domain.SigningDevice) domain.SigningDeviceFactory {
	mockFactory := mockDomain.NewMockSigningDeviceFactory(t)
	mockFactory.EXPECT().New(algo, label).Return(mockDevice, nil)
	return mockFactory
}

func TestCreateDevice(t *testing.T) {
	label := "label"
	testCreateDevice(t, &label)
}

func TestCreateDeviceNilLabel(t *testing.T) {
	var label *string = nil
	testCreateDevice(t, label)
}

func testCreateDevice(t *testing.T, label *string) {
	id := uuid.New()
	counter := 0
	signature := "signature"
	algo := string(signingapi.DeviceRequestSignatureAlgorithmRSA)

	pub := "pub"
	priv := "priv"

	mockDevice := setupMockDevice(t, id, algo, pub, priv, signature, counter, label)
	mockFactory := setupMockFactory(t, algo, label, mockDevice)
	mockStorage := mockPersistence.NewMockStorage(t)
	mockStorage.EXPECT().Add(mockDevice).Return(nil)

	dh := NewDeviceHandler(mockStorage, mockFactory)

	var reqLabel signingapi.OptNilString
	if label == nil {
		reqLabel = signingapi.OptNilString{}
	} else {
		reqLabel = signingapi.NewOptNilString(*label)
	}

	res, err := dh.CreateDevice(context.TODO(), &signingapi.DeviceRequest{
		SignatureAlgorithm: signingapi.DeviceRequestSignatureAlgorithmRSA,
		Label:              reqLabel,
	})

	assert.Nil(t, err)

	assert.Equal(t, id, res.ID)
	assert.Equal(t, pub, res.PublicKey)
	assert.Equal(t, algo, string(res.SignatureAlgorithm))
	assert.Equal(t, counter, res.Counter)
	assert.Equal(t, signature, res.LastSignature)

	resLabel, labelOK := res.Label.Get()

	if label == nil {
		assert.False(t, labelOK)
	} else {
		assert.Equal(t, *label, resLabel)
	}

}

func TestWrongAlgo(t *testing.T) {
	mockFactory := mockDomain.NewMockSigningDeviceFactory(t)
	mockStorage := mockPersistence.NewMockStorage(t)

	dh := NewDeviceHandler(mockStorage, mockFactory)

	sigalg := signingapi.DeviceRequestSignatureAlgorithm("FAKE")

	var label *string = nil

	mockFactory.EXPECT().New(string(sigalg), label).Return(nil, errors.New("Invalid algorithm"))

	_, err := dh.CreateDevice(context.TODO(), &signingapi.DeviceRequest{
		SignatureAlgorithm: sigalg,
		Label:              signingapi.OptNilString{},
	})

	assert.Error(t, err)

}
