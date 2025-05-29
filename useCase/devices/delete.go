package usecase_devices

import (
	"errors"
	"test-devices-api/database"
	domain_devices "test-devices-api/domain"
)

type DeviceDeleteErrors string

const (
	ErrorDeviceInUseCantDelete string = "the device is in use and can't be deleted"
)

func DeleteDeviceById(deviceId string) error {

	if len(deviceId) == 0 || len(deviceId) > domain_devices.MaxLenID {
		return errors.New(string(ErrorInvalidDeviceId))
	}

	deviceState, errorDeviceState := database.GetDeviceByStateID(deviceId)

	if errorDeviceState != nil {
		return errorDeviceState
	}

	if deviceState == domain_devices.DeviceStatusInUse {
		return errors.New(string(ErrorDeviceInUseCantDelete))
	}

	errorDelete := database.DeleteDeviceByID(deviceId)

	if errorDelete != nil {
		return errorDelete
	}

	return nil

}
