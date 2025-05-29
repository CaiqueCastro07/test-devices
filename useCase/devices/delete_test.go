package usecase_devices

import (
	"errors"
	"strconv"
	"test-devices-api/database"
	domain_devices "test-devices-api/domain"
	"testing"
)

func TestDeleteDeviceByID(t *testing.T) {

	c1 := domain_devices.Devices{Name: "testea", Brand: "testeb"}

	idCreation, errCreation := Create(c1)

	if errCreation != nil {
		t.Errorf("got %s, expected %s", errCreation, "nil")
		return
	}

	var testList = []struct {
		input         string
		expectedError error
	}{
		{
			input: "fja09faj0affafffffffffffffffffffffffffffffffffffffffffffffdddddddddddddd", expectedError: errors.New(string(ErrorInvalidDeviceId)),
		},
		{
			input: "", expectedError: errors.New(string(ErrorInvalidDeviceId)),
		},
		{
			input: idCreation, expectedError: nil,
		},
	}

	for i, e := range testList {

		errorDelete := DeleteDeviceById(e.input)

		if errorDelete != nil && e.expectedError != nil && e.expectedError.Error() != errorDelete.Error() {
			t.Errorf("test idx %s, got %s, expected %s", strconv.Itoa(i), errorDelete, e.expectedError)
			return
		}

	}

	device, errorGetDeviceById := database.GetDeviceByID(idCreation)

	if errorGetDeviceById != nil {
		t.Errorf("got %s, expected %s", errCreation, "nil")
		return
	}

	if len(device.Brand) > 0 {
		t.Errorf("got %s, expected %s", "existing struct", "empty struct")
		return
	}

}
