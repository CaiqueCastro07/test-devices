package database

import (
	"errors"
	"log"
	"os"
	"strconv"
	app_config "test-devices-api/config"
	domain_devices "test-devices-api/domain"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

	Connect()

}

func teardown() {
	DeleteAllDevices()
}

func TestDefaultTimeForDB(t *testing.T) {

	var now = time.Now().UTC()

	var res = defaultTimeForDB()

	if res.Sub(now).Round(0).Seconds() > 0 {
		t.Errorf("got %s, expected %s", res.Format(time.RFC3339), now.Format(time.RFC3339))
	}

}

func TestInsertDevice(t *testing.T) {

	creationId, errorInsert := InsertDevice(domain_devices.Devices{
		Brand: "test1",
		Name:  "test2",
	})

	if errorInsert != nil {
		t.Errorf("got %s, expected %s", errorInsert, "nil")
		return
	}

	if len(creationId) == 0 {
		t.Errorf("got %s, expected %s", creationId, "ID string with greater length")
		return
	}

	_, erroId := primitive.ObjectIDFromHex(creationId)

	if erroId != nil {
		t.Errorf("got %s, expected %s", erroId, "nil, for error converting HEX ID to Object ID")
		return
	}

	deviceCreated, errorDeviceCreated := GetDeviceByID(creationId)

	if errorDeviceCreated != nil {
		t.Errorf("got %s, expected %s", errorDeviceCreated, "nil")
		return
	}

	if len(deviceCreated.Brand) == 0 {
		t.Errorf("got %s, expected %s", "empty Device struct", "complete Device struct")
		return
	}

	if deviceCreated.CreationTime.IsZero() {
		t.Errorf("got %s, expected %s", "zero Device CreationTime", "UTC value for CreationTime")
		return
	}

	if deviceCreated.ModificationTime.IsZero() {
		t.Errorf("got %s, expected %s", "zero Device ModificationTime", "UTC value for ModificationTime")
		return
	}

	if deviceCreated.State != domain_devices.DefaultInitState {
		t.Errorf("got %s, expected %s", deviceCreated.State, domain_devices.DefaultInitState)
		return
	}

}

func TestUpdateDeviceByID(t *testing.T) {

	creationTime := defaultTimeForDB()

	creationId, errorInsert := InsertDevice(domain_devices.Devices{
		Brand: "test1",
		Name:  "test2",
	})

	if errorInsert != nil {
		t.Errorf("got %s, expected %s", errorInsert, "nil")
		return
	}

	var testList = []struct {
		input1      string
		input2      domain_devices.Devices
		expectedRes error
	}{
		{
			input1: creationId,
			input2: domain_devices.Devices{
				Name: "test4",
			}, expectedRes: nil,
		},
		{
			input1:      creationId,
			input2:      domain_devices.Devices{},
			expectedRes: errors.New(string(ErrorEmptyUpdateMap)),
		},
		{
			input1:      "iOAJDOajdioa",
			input2:      domain_devices.Devices{Brand: "aa"},
			expectedRes: primitive.ErrInvalidHex,
		},
	}

	lastModificationTime := defaultTimeForDB()

	for i, e := range testList {

		lastModificationTime = defaultTimeForDB()

		errUpdate := UpdateDeviceByID(e.input1, e.input2)

		if errUpdate != nil && e.expectedRes != nil && e.expectedRes.Error() != errUpdate.Error() {
			t.Errorf("test idx %s, got %s, expected %s", strconv.Itoa(i), errUpdate, e.expectedRes)
			return
		}

	}

	device, errorDevice := GetDeviceByID(creationId)

	if errorDevice != nil || len(device.Brand) == 0 {
		t.Errorf("got %s, expected %s", "error fetching updated device", "updated device data")
		return
	}

	creationTimeAfter := device.CreationTime

	if creationTimeAfter.Sub(creationTime).Round(0).Seconds() == 0 {
		t.Errorf("got %s, expected %s", "modified device's creation time", "no modification on device's creation time")
		return
	}

	if lastModificationTime == device.ModificationTime {
		t.Errorf("got %s, expected %s", "unmodified device's modification time", "modified device's modification time")
		return
	}

}

func TestGetDeviceById(t *testing.T) {

	creationId, errorInsert := InsertDevice(domain_devices.Devices{
		Brand: "test1",
		Name:  "test2",
	})

	if errorInsert != nil {
		t.Errorf("got %s, expected %s", errorInsert, "nil")
		return
	}

	deviceData, errorGetDevice := GetDeviceByID(creationId)

	if errorGetDevice != nil {
		t.Errorf("got %s, expected %s", errorGetDevice, "nil")
		return
	}

	if len(deviceData.Brand) == 0 {
		t.Errorf("got %s, expected %s", "no brand", "test1")
		return
	}

	if len(deviceData.Name) == 0 {
		t.Errorf("got %s, expected %s", "no name", "test2")
		return
	}

	errorDeleting := DeleteDeviceByID(creationId)

	if errorDeleting != nil {
		t.Errorf("got %s, expected %s", "no name", "test2")
		return
	}

	deviceData2, errorGetDevice2 := GetDeviceByID(creationId)

	if errorGetDevice2 != nil {
		t.Errorf("got %s, expected %s", errorGetDevice2, "nil")
		return
	}

	if len(deviceData2.Brand) != 0 {
		t.Errorf("got %s, expected %s", deviceData2.Brand, "empty struct")
		return
	}

}

func TestGetAllDevicesByField(t *testing.T) {

	c1 := domain_devices.Devices{
		Brand: "apple",
		Name:  "iphone 10",
	}

	c2 := domain_devices.Devices{
		Brand: "apple",
		Name:  "iphone 11",
	}

	c3 := domain_devices.Devices{
		Brand: "samsung",
		Name:  "galaxy 11",
	}

	c4 := domain_devices.Devices{
		Brand: "samsung",
		Name:  "galaxy 11",
	}

	toCreate := []domain_devices.Devices{c1, c2, c3, c4}

	for _, e := range toCreate {

		_, databaseError := InsertDevice(e)

		if databaseError != nil {
			t.Errorf("got %s, expected %s", databaseError, "nil")
			return
		}

	}

	list1, dbError := GetAllDevicesByField(domain_devices.Devices{
		Brand: c1.Brand,
	})

	if dbError != nil {
		t.Errorf("got %s, expected %s", dbError, "nil")
		return
	}

	if len(*list1) != 2 {
		t.Errorf("got %s, expected %s", "list len "+strconv.Itoa(len(*list1)), "2")
		return
	}

	list2, dbError := GetAllDevicesByField(domain_devices.Devices{
		Name: c3.Name,
	})

	if dbError != nil {
		t.Errorf("got %s, expected %s", dbError, "nil")
		return
	}

	if len(*list2) != 2 {
		t.Errorf("got %s, expected %s", "list len "+strconv.Itoa(len(*list2)), "2")
		return
	}

	list3, dbError := GetAllDevicesByField(domain_devices.Devices{})

	if dbError != nil {
		t.Errorf("got %s, expected %s", dbError, "nil")
		return
	}

	if len(*list3) < 4 {
		t.Errorf("got %s, expected %s", "list len "+strconv.Itoa(len(*list3)), "GTE than 4 brings every document")
		return
	}

}

func TestDeleteDeviceByid(t *testing.T) {

	dbError := DeleteAllDevices()

	if dbError != nil {
		t.Errorf("got %s, expected %s", dbError, "nil")
		return
	}

	time.Sleep(3 * time.Second)

	list, dbError := GetAllDevicesByField(domain_devices.Devices{})

	if dbError != nil {
		t.Errorf("got %s, expected %s", dbError, "nil")
		return
	}

	if len(*list) != 0 {
		t.Errorf("got %s, expected %s", strconv.Itoa(len(*list)), "0")
	}

}
