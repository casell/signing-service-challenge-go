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

func TestSignTransactionGetError(t *testing.T) {
	id := uuid.New()
	dataToBeSigned := "data"

	mockFactory := mockDomain.NewMockSigningDeviceFactory(t)

	mockStorage := mockPersistence.NewMockStorage(t)
	getErr := errors.New("Get Error")
	mockStorage.EXPECT().Get(id).Return(nil, getErr)

	dh := NewDeviceHandler(mockStorage, mockFactory)

	params := signingapi.SignTransactionParams{
		Deviceid: id,
	}

	req := &signingapi.SignatureRequest{
		DataToBeSigned: dataToBeSigned,
	}

	res, err := dh.SignTransaction(context.TODO(), req, params)

	assert.Nil(t, res)
	if assert.Error(t, err) {
		assert.Equal(t, getErr, err)
	}
}

func TestSignTransactionNoDevice(t *testing.T) {
	id := uuid.New()
	dataToBeSigned := "data"

	mockFactory := mockDomain.NewMockSigningDeviceFactory(t)

	mockStorage := mockPersistence.NewMockStorage(t)
	mockStorage.EXPECT().Get(id).Return(nil, nil)

	dh := NewDeviceHandler(mockStorage, mockFactory)

	params := signingapi.SignTransactionParams{
		Deviceid: id,
	}

	req := &signingapi.SignatureRequest{
		DataToBeSigned: dataToBeSigned,
	}

	res, err := dh.SignTransaction(context.TODO(), req, params)

	assert.Nil(t, res)
	if assert.Error(t, err) {
		assert.IsType(t, errDeviceNotFound{}, err)
	}
}

func TestSignTransactionErrorSigning(t *testing.T) {
	id := uuid.New()
	dataToBeSigned := "data"

	mockFactory := mockDomain.NewMockSigningDeviceFactory(t)

	mockStorage := mockPersistence.NewMockStorage(t)

	mockDevice := mockDomain.NewMockSigningDevice(t)
	// mockDevice.EXPECT().ID().Return(id)

	signErr := errors.New("Sign error")

	mockDevice.EXPECT().Sign(dataToBeSigned).Return("", "", signErr)

	mockStorage.EXPECT().Get(id).Return(mockDevice, nil)

	dh := NewDeviceHandler(mockStorage, mockFactory)

	params := signingapi.SignTransactionParams{
		Deviceid: id,
	}

	req := &signingapi.SignatureRequest{
		DataToBeSigned: dataToBeSigned,
	}

	res, err := dh.SignTransaction(context.TODO(), req, params)

	assert.Nil(t, res)
	if assert.Error(t, err) {
		assert.Equal(t, signErr, err)
	}
}

func TestSignTransactionErrorPut(t *testing.T) {
	id := uuid.New()
	dataToBeSigned := "data"

	mockFactory := mockDomain.NewMockSigningDeviceFactory(t)

	mockStorage := mockPersistence.NewMockStorage(t)

	mockDevice := mockDomain.NewMockSigningDevice(t)

	signature := "Signature"
	extData := "Ext Data"

	mockDevice.EXPECT().Sign(dataToBeSigned).Return(signature, extData, nil)

	mockStorage.EXPECT().Get(id).Return(mockDevice, nil)

	putErr := errors.New("Put error")
	mockStorage.EXPECT().Put(mockDevice).Return(putErr)

	dh := NewDeviceHandler(mockStorage, mockFactory)

	params := signingapi.SignTransactionParams{
		Deviceid: id,
	}

	req := &signingapi.SignatureRequest{
		DataToBeSigned: dataToBeSigned,
	}

	res, err := dh.SignTransaction(context.TODO(), req, params)

	assert.Nil(t, res)
	if assert.Error(t, err) {
		assert.Equal(t, putErr, err)
	}
}

func TestSignTransaction(t *testing.T) {
	id := uuid.New()
	dataToBeSigned := "data"

	mockFactory := mockDomain.NewMockSigningDeviceFactory(t)

	mockStorage := mockPersistence.NewMockStorage(t)

	mockDevice := mockDomain.NewMockSigningDevice(t)

	signature := "Signature"
	extData := "Ext Data"

	mockDevice.EXPECT().Sign(dataToBeSigned).Return(signature, extData, nil)

	mockStorage.EXPECT().Get(id).Return(mockDevice, nil)

	mockStorage.EXPECT().Put(mockDevice).Return(nil)

	dh := NewDeviceHandler(mockStorage, mockFactory)

	params := signingapi.SignTransactionParams{
		Deviceid: id,
	}

	req := &signingapi.SignatureRequest{
		DataToBeSigned: dataToBeSigned,
	}

	res, err := dh.SignTransaction(context.TODO(), req, params)

	assert.Nil(t, err)
	assert.Equal(t, signature, res.GetSignature())
	assert.Equal(t, extData, res.GetSignedData())
}
