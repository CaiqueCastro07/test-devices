package routes

import (
	"log"
	"net/http"
	"net/http/httptest"
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

var han *http.ServeMux

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

func TestHandlerGetAllDevicesReq(t *testing.T) {

	handler := http.NewServeMux()
	handler.HandleFunc(DevicesPath, GetAllDevicesReq)

	// Set up test database
	c1 := domain_devices.Devices{
		Name:  "testhh",
		Brand: "testkk",
	}

	c2 := domain_devices.Devices{
		Name:  "testzz",
		Brand: "testyy",
	}

	var dList = []domain_devices.Devices{c1, c2}

	for _, e := range dList {

		_, errorInsertion := database.InsertDevice(e)

		if errorInsertion != nil {
			t.Errorf("got %s, expected %s", errorInsertion, "nil")
			return
		}

	}

	var reqLists = []struct {
		query          string
		expectedStatus int
	}{
		{query: "?brand=noo", expectedStatus: http.StatusOK},
		{query: "?brand=" + c1.Brand, expectedStatus: http.StatusOK},
		{query: "", expectedStatus: http.StatusOK},
	}

	var reqType = "GET"

	for _, e := range reqLists {

		req, err := http.NewRequest(reqType, DevicesPath+e.query, nil)

		if err != nil {
			t.Errorf("failed to create request client")
		}

		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		if res.Result().StatusCode != e.expectedStatus {
			t.Errorf("got %s, expected %s", strconv.Itoa(res.Result().StatusCode), strconv.Itoa(e.expectedStatus))
		}
		// TODO
		// validate response data

	}

}

func TestHandlerGetDeviceByIDReq(t *testing.T) {

	handler := http.NewServeMux()
	handler.HandleFunc(DevicesPath, GetAllDevicesReq)
	// TODO
}

func TestHandlerCreateDeviceReq(t *testing.T) {

	handler := http.NewServeMux()
	handler.HandleFunc(DevicesPath, CreateDeviceReq)
	// TODO
}

func TestHandlerUpdateDeviceReq(t *testing.T) {

	handler := http.NewServeMux()
	handler.HandleFunc(DevicesPath, UpdateDeviceReq)
	// TODO

}

func TestHandlerDeleteDeviceByIDReq(t *testing.T) {

	handler := http.NewServeMux()
	handler.HandleFunc(DevicesPath, DeleteDeviceByIDReq)
	// TODO

}
