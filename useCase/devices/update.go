package usecase_devices

import (
	"errors"
	"strings"
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

	errorValidacao := validateDeviceDataForUpdate(&deviceData)

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

func validateDeviceDataForUpdate(deviceData *domain_devices.Devices) error {

	if len(deviceData.Brand) > 0 {

		if brandLen := len(deviceData.Brand); brandLen > domain_devices.MaxLenBrand {
			return errors.New(string(ErrorInvalidBrand))
		}

		deviceData.Brand = strings.ToLower(strings.Trim(deviceData.Brand, " "))

	}

	if len(deviceData.Name) > 0 {

		if nameLen := len(deviceData.Name); nameLen > domain_devices.MaxLenName {
			return errors.New(string(ErrorInvalidName))
		}

		deviceData.Name = strings.ToLower(strings.Trim(deviceData.Name, " "))

	}

	if len(deviceData.State) > 0 {

		errorValidatingState := validateState(deviceData.State)

		if errorValidatingState != nil {
			return errorValidatingState
		}

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
