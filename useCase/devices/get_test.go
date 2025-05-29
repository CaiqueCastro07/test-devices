package usecase_devices

import (
	"errors"
	"strconv"
	"test-devices-api/database"
	domain_devices "test-devices-api/domain"
	"testing"
	"time"
)

func TestGetDeviceById(t *testing.T) {

	c1 := domain_devices.Devices{Name: "testea", Brand: "testeb"}

	idCreation, errCreation := Create(c1)

	if errCreation != nil {
		t.Errorf("got %s, expected %s", errCreation, "nil")
		return
	}

	deviceData, errorGet := GetDeviceById(idCreation)

	if errorGet != nil {
		t.Errorf("got %s, expected %s", errorGet, "nil")
		return
	}

	if deviceData.Brand != c1.Brand {
		t.Errorf("got %s, expected %s", "brand "+deviceData.Brand, c1.Brand)
		return
	}

	if deviceData.Name != c1.Name {
		t.Errorf("got %s, expected %s", "name "+deviceData.Name, c1.Name)
		return
	}

	errorDeleting := database.DeleteDeviceByID(idCreation)

	if errorDeleting != nil {
		t.Errorf("got %s, expected %s", errorDeleting, "nil")
		return
	}

	time.Sleep(1 * time.Second)

	deviceData2, errorGet := GetDeviceById(idCreation)

	if errorGet != nil {
		t.Errorf("got %s, expected %s", errorGet, "nil")
		return
	}

	if len(deviceData2.Brand) != 0 {
		t.Errorf("got %s, expected %s", "existing struct", "empty struct")
		return
	}

	_, errorGet2 := GetDeviceById("")

	if errorGet2 != nil && errorGet2.Error() != ErrorInvalidDeviceId {
		t.Errorf("got %s, expected %s", errorGet2, ErrorInvalidDeviceId)
		return
	}

	_, errorGet3 := GetDeviceById("fJ8FJ89AF89AJ9F8AJ89FJ89AJF89AJF98AJ9F8AJ89FJA98JA89GJA89GJA98GJA98GJ98AJD89AJF98AJ9FAJ89FJA8989FA9J")

	if errorGet3 != nil && errorGet3.Error() != ErrorInvalidDeviceId {
		t.Errorf("got %s, expected %s", errorGet3, ErrorInvalidDeviceId)
		return
	}

}

func TestGetAllDevicesWithFilters(t *testing.T) {

	var testList = []struct {
		input         domain_devices.Devices
		expectedError error
	}{
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
				State: "AvaLaible",
			},
			expectedError: errors.New(string(ErrorInvalidState)),
		},
		{
			input:         domain_devices.Devices{},
			expectedError: nil,
		},
		{
			input:         domain_devices.Devices{State: domain_devices.DeviceStatusInUse},
			expectedError: nil,
		},
	}

	for i, e := range testList {

		_, errorGet := GetAllDevicesWithFilters(e.input)

		if errorGet != nil && e.expectedError != nil && e.expectedError.Error() != errorGet.Error() {
			t.Errorf("test idx %s, got %s, expected %s", strconv.Itoa(i), errorGet, e.expectedError)
			return
		}

	}

	c1 := domain_devices.Devices{
		Name:  "testeC",
		Brand: "testec",
	}

	_, errorCreation := Create(c1)

	if errorCreation != nil {
		t.Errorf("got %s, expected %s", errorCreation, "nil")
		return
	}

	time.Sleep(1 * time.Second)

	list, errorGet := GetAllDevicesWithFilters(c1)

	if errorGet != nil {
		t.Errorf("got %s, expected %s", errorGet, "nil")
		return
	}

	if len(*list) != 1 {
		t.Errorf("got %s, expected %s", "list with "+strconv.Itoa(len(*list)), "1")
		return
	}

}

func TestValidateDeviceDataForFilter(t *testing.T) {

	var testList = []struct {
		input         domain_devices.Devices
		expectedError error
	}{
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
				State: "AvaLaible",
			},
			expectedError: errors.New(string(ErrorInvalidState)),
		},
		{
			input:         domain_devices.Devices{},
			expectedError: nil,
		},
		{
			input:         domain_devices.Devices{State: domain_devices.DeviceStatusInUse},
			expectedError: nil,
		},
	}

	for i, e := range testList {

		errValidation := validateDeviceDataForFilter(&e.input)

		if errValidation != nil && e.expectedError != nil && e.expectedError.Error() != errValidation.Error() {
			t.Errorf("test idx %s, got %s, expected %s", strconv.Itoa(i), errValidation, e.expectedError)
			return
		}

	}

	var testList2 = []struct {
		input       domain_devices.Devices
		expectedRes domain_devices.Devices
	}{
		{
			input: domain_devices.Devices{
				Name:  "TESTE4 ",
				Brand: "   ATT",
			}, expectedRes: domain_devices.Devices{
				Name:  "teste4",
				Brand: "att",
			},
		},
	}

	for i, e := range testList2 {

		errValidation := validateDeviceDataForFilter(&e.input)

		if errValidation != nil {
			t.Errorf("test idx %s, got %s, expected %s", strconv.Itoa(i), errValidation, "")
			return
		}

		if e.input.Brand != e.expectedRes.Brand {
			t.Errorf("test idx %s, got %s, expected %s", strconv.Itoa(i), e.input.Brand, e.expectedRes.Brand)
			return
		}

		if e.input.Name != e.expectedRes.Name {
			t.Errorf("test idx %s, got %s, expected %s", strconv.Itoa(i), e.input.Name, e.expectedRes.Name)
			return
		}

	}

}
