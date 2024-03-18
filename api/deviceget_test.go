package api

import (
	"context"
	"errors"
	"testing"

	"github.com/casell/signing-service-challenge/generated/signingapi"
	mockDomain "github.com/casell/signing-service-challenge/mocks/github.com/casell/signing-service-challenge/domain"
	mockPersistence "github.com/casell/signing-service-challenge/mocks/github.com/casell/signing-service-challenge/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetDevice(t *testing.T) {
	id := uuid.New()
	counter := 0
	signature := "signature"
	algo := string(signingapi.DeviceRequestSignatureAlgorithmRSA)

	pub := "pub"
	priv := "priv"

	label := "label"

	mockDevice := setupMockDevice(t, id, algo, pub, priv, signature, counter, &label)

	mockFactory := mockDomain.NewMockSigningDeviceFactory(t)

	mockStorage := mockPersistence.NewMockStorage(t)

	mockStorage.EXPECT().Get(id).Return(mockDevice, nil)

	dh := NewDeviceHandler(mockStorage, mockFactory)

	params := signingapi.GetDeviceParams{
		Deviceid: id,
	}

	res, err := dh.GetDevice(context.TODO(), params)

	assert.Nil(t, err)

	assert.Equal(t, id, res.ID)
	assert.Equal(t, pub, res.PublicKey)
	assert.Equal(t, algo, string(res.SignatureAlgorithm))
	assert.Equal(t, counter, res.Counter)
	assert.Equal(t, signature, res.LastSignature)

	resLabel, labelOK := res.Label.Get()

	assert.True(t, labelOK)
	assert.Equal(t, label, resLabel)
}

func TestGetDeviceError(t *testing.T) {
	id := uuid.New()

	mockFactory := mockDomain.NewMockSigningDeviceFactory(t)

	mockStorage := mockPersistence.NewMockStorage(t)
	getErr := errors.New("Get Error")
	mockStorage.EXPECT().Get(id).Return(nil, getErr)

	dh := NewDeviceHandler(mockStorage, mockFactory)

	params := signingapi.GetDeviceParams{
		Deviceid: id,
	}

	res, err := dh.GetDevice(context.TODO(), params)

	assert.Nil(t, res)
	if assert.Error(t, err) {
		assert.Equal(t, getErr, err)
	}
}

func TestGetDeviceNotFound(t *testing.T) {
	id := uuid.New()

	mockFactory := mockDomain.NewMockSigningDeviceFactory(t)

	mockStorage := mockPersistence.NewMockStorage(t)
	mockStorage.EXPECT().Get(id).Return(nil, nil)

	dh := NewDeviceHandler(mockStorage, mockFactory)

	params := signingapi.GetDeviceParams{
		Deviceid: id,
	}

	res, err := dh.GetDevice(context.TODO(), params)

	assert.Nil(t, res)
	if assert.Error(t, err) {
		assert.IsType(t, errDeviceNotFound{}, err)
	}
}
