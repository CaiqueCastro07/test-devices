package usecase_devices

import (
	"errors"
	"strings"
	"test-devices-api/database"
	domain_devices "test-devices-api/domain"
)

type DeviceGetErrors string

const ()

func GetDeviceById(deviceId string) (*domain_devices.Devices, error) {

	if len(deviceId) == 0 || len(deviceId) > domain_devices.MaxLenID {
		return &domain_devices.Devices{}, errors.New(string(ErrorInvalidDeviceId))
	}

	device, errorGet := database.GetDeviceByID(deviceId)

	if errorGet != nil {
		return &domain_devices.Devices{}, errorGet
	}

	return device, nil

}

func GetAllDevicesWithFilters(filters domain_devices.Devices) (*[]domain_devices.Devices, error) {

	errorValidatingFilters := validateDeviceDataForFilter(&filters)

	if errorValidatingFilters != nil {
		return &[]domain_devices.Devices{}, errorValidatingFilters
	}

	lista, errorGet := database.GetAllDevicesByField(filters)

	if errorGet != nil {
		return lista, errorGet
	}

	return lista, nil

}

func validateDeviceDataForFilter(deviceData *domain_devices.Devices) error {

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
