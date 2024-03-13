package api

import (
	"context"
	"errors"
	"testing"

	"github.com/casell/signing-service-challenge/domain"
	mockDomain "github.com/casell/signing-service-challenge/mocks/github.com/casell/signing-service-challenge/domain"
	mockPersistence "github.com/casell/signing-service-challenge/mocks/github.com/casell/signing-service-challenge/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListDevice(t *testing.T) {
	mockDevice := mockDomain.NewMockSigningDevice(t)
	id := uuid.New()
	label := "label"
	mockDevice.EXPECT().ID().Return(id)
	mockDevice.EXPECT().Label().Return(&label)

	mockFactory := mockDomain.NewMockSigningDeviceFactory(t)

	mockStorage := mockPersistence.NewMockStorage(t)
	mockStorage.EXPECT().List().Return([]domain.SigningDevice{mockDevice}, nil)

	dh := NewDeviceHandler(mockStorage, mockFactory)

	res, err := dh.ListDevices(context.TODO())

	assert.Nil(t, err)
	if assert.Len(t, res, 1) {
		assert.Equal(t, id, res[0].ID)
		lab, ok := res[0].GetLabel().Get()
		assert.True(t, ok)
		assert.Equal(t, label, lab)
	}
}

func TestListDeviceEmpty(t *testing.T) {
	mockFactory := mockDomain.NewMockSigningDeviceFactory(t)

	mockStorage := mockPersistence.NewMockStorage(t)
	mockStorage.EXPECT().List().Return([]domain.SigningDevice{}, nil)

	dh := NewDeviceHandler(mockStorage, mockFactory)

	res, err := dh.ListDevices(context.TODO())

	assert.Nil(t, err)
	assert.Len(t, res, 0)
}

func TestListDeviceError(t *testing.T) {
	listErr := errors.New("List error")

	mockFactory := mockDomain.NewMockSigningDeviceFactory(t)

	mockStorage := mockPersistence.NewMockStorage(t)

	mockStorage.EXPECT().List().Return(nil, listErr)

	dh := NewDeviceHandler(mockStorage, mockFactory)

	res, err := dh.ListDevices(context.TODO())

	assert.Nil(t, res)
	if assert.Error(t, err) {
		assert.Equal(t, listErr, err)
	}
}
