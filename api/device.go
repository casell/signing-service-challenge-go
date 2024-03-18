package api

import (
	"context"
	"net/http"

	"github.com/casell/signing-service-challenge/domain"
	"github.com/casell/signing-service-challenge/generated/signingapi"
	"github.com/casell/signing-service-challenge/persistence"
)

// DeviceHandler represents the HTTP Handler to reply to signing api requests.
type DeviceHandler struct {
	store         persistence.Storage
	devicefactory domain.SigningDeviceFactory
}

// NewDeviceHandler creates a device handler backed by the Storage store.
func NewDeviceHandler(store persistence.Storage, devicefactory domain.SigningDeviceFactory) *DeviceHandler {
	return &DeviceHandler{
		store:         store,
		devicefactory: devicefactory,
	}
}

// CreateDevice handles device creation requests.
func (h *DeviceHandler) CreateDevice(ctx context.Context, req *signingapi.DeviceRequest) (*signingapi.DeviceResponse, error) {

	origlabel, ok := req.GetLabel().Get()
	var label *string

	if !ok {
		label = nil
	} else {
		label = &origlabel
	}

	device, err := h.devicefactory.New(string(req.GetSignatureAlgorithm()), label)
	if err != nil {
		return nil, err
	}
	if err := h.store.Add(device); err != nil {
		return nil, err
	}

	return convertToApiResponse(device)
}

// SignTransaction handles signing requests.
func (h *DeviceHandler) SignTransaction(ctx context.Context, req *signingapi.SignatureRequest, params signingapi.SignTransactionParams) (*signingapi.SignatureResponse, error) {

	device, err := h.store.Get(params.Deviceid)
	if err != nil {
		return nil, err
	}

	if device == nil {
		return nil, errDeviceNotFound{params.Deviceid.String()}
	}

	signature, extendedDataToBeSigned, err := device.Sign(req.DataToBeSigned)
	if err != nil {
		return nil, err
	}

	if err := h.store.Put(device); err != nil {
		return nil, err
	}

	return &signingapi.SignatureResponse{
		Signature:  signature,
		SignedData: extendedDataToBeSigned,
	}, nil
}

// ListDevices handles device list requests.
func (h *DeviceHandler) ListDevices(ctx context.Context) ([]signingapi.DeviceSummary, error) {
	devices, err := h.store.List()
	if err != nil {
		return nil, err
	}

	devicesummaries := make([]signingapi.DeviceSummary, len(devices))
	for i, v := range devices {
		optlabel := signingapi.OptString{}
		if label := v.Label(); label != nil {
			optlabel.SetTo(*label)
		}
		devicesummaries[i] = signingapi.DeviceSummary{
			ID:    v.ID(),
			Label: optlabel,
		}
	}
	return devicesummaries, nil
}

// GetDevice handles device retrieval requests
func (h *DeviceHandler) GetDevice(ctx context.Context, params signingapi.GetDeviceParams) (*signingapi.DeviceResponse, error) {
	device, err := h.store.Get(params.Deviceid)
	if err != nil {
		return nil, err
	}

	if device == nil {
		return nil, errDeviceNotFound{params.Deviceid.String()}
	}

	return convertToApiResponse(device)
}

func convertToApiResponse(device domain.SigningDevice) (*signingapi.DeviceResponse, error) {

	counter, lastSignature := device.CounterAndLastSignature()
	pub, _, err := device.KeyPair().Marshal()
	if err != nil {
		return nil, err
	}

	var sigalg signingapi.DeviceResponseSignatureAlgorithm
	if err := sigalg.UnmarshalText([]byte(device.SignatureAlgorithm())); err != nil {
		return nil, err
	}

	optlabel := signingapi.OptString{}

	if label := device.Label(); label != nil {
		optlabel.SetTo(*label)
	}

	return &signingapi.DeviceResponse{
		ID:                 device.ID(),
		SignatureAlgorithm: sigalg,
		Label:              optlabel,
		Counter:            int(counter),
		LastSignature:      lastSignature,
		PublicKey:          string(pub),
	}, nil
}

// NewError converts errors to an http structure response
func (h *DeviceHandler) NewError(ctx context.Context, err error) *signingapi.ErrorResponseStatusCode {
	switch err.(type) {
	case domain.ErrInvalidAlgorithm:
		return &signingapi.ErrorResponseStatusCode{
			StatusCode: http.StatusBadRequest,
			Response: signingapi.ErrorResponse{
				Errors: []string{err.Error()},
			},
		}
	case errDeviceNotFound:
		return &signingapi.ErrorResponseStatusCode{
			StatusCode: http.StatusNotFound,
			Response: signingapi.ErrorResponse{
				Errors: []string{err.Error()},
			},
		}
	default:
		return &signingapi.ErrorResponseStatusCode{
			StatusCode: http.StatusInternalServerError,
			Response: signingapi.ErrorResponse{
				Errors: []string{err.Error()},
			},
		}
	}
}
