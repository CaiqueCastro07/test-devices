package usecase_devices

import (
	"errors"
	"strings"
	"test-devices-api/database"
	domain_devices "test-devices-api/domain"
)

type DeviceErrors string

type DeviceCreationErrors string

const (
	ErrorInvalidBrand DeviceErrors = "invalid brand"
	ErrorInvalidName  DeviceErrors = "invalid name"
)

var ErrorInvalidState string = "invalid state, expected " + strings.Join(domain_devices.DeviceStatesList[:], " or ")

func Create(deviceData domain_devices.Devices) (string, error) {

	errorValidacao := validateDeviceDataForCreation(&deviceData)

	if errorValidacao != nil {
		return "", errorValidacao
	}

	creationId, errorCreation := database.InsertDevice(deviceData)

	if errorCreation != nil {
		return "", errorCreation
	}

	return creationId, nil

}

func validateDeviceDataForCreation(deviceData *domain_devices.Devices) error {

	if brandLen := len(deviceData.Brand); brandLen == 0 || brandLen > 30 {
		return errors.New(string(ErrorInvalidBrand))
	}

	deviceData.Brand = strings.ToLower(strings.Trim(deviceData.Brand, " "))

	if nameLen := len(deviceData.Name); nameLen == 0 || nameLen > 30 {
		return errors.New(string(ErrorInvalidName))
	}

	deviceData.Name = strings.ToLower(strings.Trim(deviceData.Name, " "))

	return nil

}
