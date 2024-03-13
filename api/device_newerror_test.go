package api

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/casell/signing-service-challenge/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewErrorInvalidAlgorithm(t *testing.T) {
	var dh *DeviceHandler

	err := domain.ErrInvalidAlgorithm{}
	errResp := dh.NewError(context.TODO(), err)
	assert.Equal(t, http.StatusBadRequest, errResp.GetStatusCode())
	if assert.NotNil(t, errResp.GetResponse()) {
		if assert.Len(t, errResp.GetResponse().Errors, 1) {
			assert.Equal(t, err.Error(), errResp.GetResponse().Errors[0])
		}
	}
}

func TestNewErrorDeviceNotFound(t *testing.T) {
	var dh *DeviceHandler

	err := errDeviceNotFound{deviceID: "42"}
	errResp := dh.NewError(context.TODO(), err)
	assert.Equal(t, http.StatusNotFound, errResp.GetStatusCode())
	if assert.NotNil(t, errResp.GetResponse()) {
		if assert.Len(t, errResp.GetResponse().Errors, 1) {
			assert.Equal(t, err.Error(), errResp.GetResponse().Errors[0])
		}
	}
}

func TestNewErrorDefault(t *testing.T) {
	var dh *DeviceHandler

	err := errors.New("Error error!")
	errResp := dh.NewError(context.TODO(), err)
	assert.Equal(t, http.StatusInternalServerError, errResp.GetStatusCode())
	if assert.NotNil(t, errResp.GetResponse()) {
		if assert.Len(t, errResp.GetResponse().Errors, 1) {
			assert.Equal(t, err.Error(), errResp.GetResponse().Errors[0])
		}
	}
}
