package usecase_devices

import (
	"errors"
	"strconv"
	"test-devices-api/database"
	domain_devices "test-devices-api/domain"
	"testing"
)

func TestUpdate(t *testing.T) {

	c1 := domain_devices.Devices{
		Name:  "testeG",
		Brand: "testeG",
	}

	creationId, errorCreation := Create(c1)

	if errorCreation != nil {
		t.Errorf("got %s, expected %s", errorCreation, "nil")
		return
	}

	var testList = []struct {
		input         domain_devices.Devices
		expectedError error
	}{
		{
			input: domain_devices.Devices{}, expectedError: errors.New(string(database.ErrorEmptyUpdateMap)),
		},
		{
			input: domain_devices.Devices{
				Name: "test4aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			}, expectedError: errors.New(string(ErrorInvalidName)),
		},
		{
			input: domain_devices.Devices{
				Brand: "0j1f0201290fj0fffffffffffffffffaaaaaaaaaaaaaaaaaaaaaaaaaafffffffffffffaaaaaaaaaaaffffffffffff",
			},
			expectedError: errors.New(string(ErrorInvalidBrand)),
		},
		{
			input: domain_devices.Devices{
				Brand: "novo",
			},
			expectedError: nil,
		},
		{
			input: domain_devices.Devices{
				State: "AvaLaible",
			},
			expectedError: errors.New(string(ErrorInvalidState)),
		},
		{
			input:         domain_devices.Devices{State: domain_devices.DeviceStatusInUse},
			expectedError: nil,
		},
		{
			input:         domain_devices.Devices{Brand: "NOVO"},
			expectedError: errors.New(ErrorDeviceInUse),
		},
	}

	for i, e := range testList {

		errorUpdate := Update(creationId, e.input)

		if errorUpdate != nil && e.expectedError != nil && e.expectedError.Error() != errorUpdate.Error() {
			t.Errorf("test idx %s, got %s, expected %s", strconv.Itoa(i), errorUpdate, e.expectedError)
			return
		}

	}

}
