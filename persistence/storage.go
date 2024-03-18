package persistence

import (
	"github.com/casell/signing-service-challenge/domain"
	"github.com/google/uuid"
)

type Storage interface {
	List() ([]domain.SigningDevice, error)
	Get(id uuid.UUID) (domain.SigningDevice, error)
	Add(x domain.SigningDevice) error
	Put(x domain.SigningDevice) error
}
