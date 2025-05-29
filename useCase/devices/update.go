package usecase_devices

import (
	"errors"
	"test-devices-api/database"
	domain_devices "test-devices-api/domain"
)

type DeviceUpdateErrors string

const (
	ErrorInvalidDeviceId string = "invalid device id"
	ErrorDeviceInUse     string = "can't update a device currently in use"
)

func Update(deviceId string, deviceData domain_devices.Devices) error {

	if len(deviceId) == 0 || len(deviceId) > domain_devices.MaxLenID {
		return errors.New(string(ErrorInvalidDeviceId))
	}

	errorValidacao := validateDeviceDataForFilter(&deviceData)

	if errorValidacao != nil {
		return errorValidacao
	}

	deviceState, errorDeviceState := database.GetDeviceByStateID(deviceId)

	if errorDeviceState != nil {
		return errorDeviceState
	}

	if deviceState == domain_devices.DeviceStatusInUse {
		return errors.New(string(ErrorDeviceInUse))
	}

	errorUpdate := database.UpdateDeviceByID(deviceId, deviceData)

	if errorUpdate != nil {
		return errorUpdate
	}

	return nil

}

func validateState(state domain_devices.DeviceState) error {

	if len(state) == 0 {
		return nil
	}

	if state != domain_devices.DeviceStatusAvalaible && state != domain_devices.DeviceStatusInUse &&
		state != domain_devices.DeviceStatusInactive {
		return errors.New(ErrorInvalidState)
	}

	return nil

}
