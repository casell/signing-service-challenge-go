package api

import (
	"fmt"
)

type errDeviceNotFound struct {
	deviceID string
}

func (e errDeviceNotFound) Error() string {
	return fmt.Sprintf("device %s not found", e.deviceID)

}
