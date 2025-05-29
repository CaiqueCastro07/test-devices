package usecase_devices

import (
	"errors"
	"log"
	"os"
	"strconv"
	app_config "test-devices-api/config"
	"test-devices-api/database"
	domain_devices "test-devices-api/domain"
	"testing"
)

func TestMain(m *testing.M) {
	mainSetup()
	exitCode := m.Run()
	teardown()
	os.Exit(exitCode)
}

func mainSetup() {

	errorSet := os.Setenv(app_config.DEFAULT_KEY_FOR_CONFIG, "TEST")

	if errorSet != nil {
		log.Fatal(errorSet)
	}

	app_config.LoadConfig()

	database.Connect()

}

func teardown() {
	database.DeleteAllDevices()
}

func TestCreate(t *testing.T) {

	var testList = []struct {
		input         domain_devices.Devices
		expectedError error
		expectedRes   string
	}{
		{
			input: domain_devices.Devices{
				Name: "test4",
			}, expectedError: errors.New(string(ErrorInvalidBrand)),
			expectedRes: "",
		},
		{
			input: domain_devices.Devices{
				Brand: "0j1f0201290fj0fffffffffffffffffaaaaaaaaaaaaaaaaaaaaaaaaaafffffffffffffffffffffffff",
				Name:  "test4",
			},
			expectedError: errors.New(string(ErrorInvalidBrand)),
			expectedRes:   "",
		},
		{
			input:         domain_devices.Devices{Name: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Brand: "test2"},
			expectedError: errors.New(string(ErrorInvalidName)),
			expectedRes:   "",
		},
		{
			input:         domain_devices.Devices{Brand: "test2"},
			expectedError: errors.New(string(ErrorInvalidName)),
			expectedRes:   "",
		},
		{
			input:         domain_devices.Devices{Brand: "test2", Name: "test1"},
			expectedError: nil,
			expectedRes:   "id",
		},
	}

	for i, e := range testList {

		idCreation, errCreation := Create(e.input)

		if errCreation != nil && e.expectedError != nil && e.expectedError.Error() != errCreation.Error() {
			t.Errorf("test idx %s, got %s, expected %s", strconv.Itoa(i), errCreation, e.expectedError)
			return
		}

		if e.expectedRes == "id" && len(idCreation) == 0 {
			t.Errorf("test idx %s, got %s, expected %s", strconv.Itoa(i), "empty string", "creation id")
			return
		}

	}

}

func TestValidateDeviceDataForCreation(t *testing.T) {

	var testList = []struct {
		input         domain_devices.Devices
		expectedError error
	}{
		{
			input: domain_devices.Devices{
				Name: "test4",
			}, expectedError: errors.New(string(ErrorInvalidBrand)),
		},
		{
			input: domain_devices.Devices{
				Brand: "0j1f0201290fj0fffffffffffffffffaaaaaaaaaaaaaaaaaaaaaaaaaafffffffffffffffffffffffff",
				Name:  "test4",
			},
			expectedError: errors.New(string(ErrorInvalidBrand)),
		},
		{
			input:         domain_devices.Devices{Name: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Brand: "test2"},
			expectedError: errors.New(string(ErrorInvalidName)),
		},
		{
			input:         domain_devices.Devices{Brand: "test2"},
			expectedError: errors.New(string(ErrorInvalidName)),
		},
		{
			input:         domain_devices.Devices{Brand: "test2", Name: "test1"},
			expectedError: nil,
		},
	}

	for i, e := range testList {

		errValidation := validateDeviceDataForCreation(&e.input)

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

		errValidation := validateDeviceDataForCreation(&e.input)

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
