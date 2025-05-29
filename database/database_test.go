package database

import (
	"fmt"
	"os"
	app_config "test-devices-api/config"
	domain_devices "test-devices-api/domain"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setup() {

	errorSet := os.Setenv(app_config.DEFAULT_KEY_FOR_CONFIG, "TEST")

	if errorSet != nil {
		fmt.Println("errorset", errorSet)
	}

	app_config.LoadConfig()

	Connect()

}

func teardown() {
	fmt.Println("Tearing down after tests")
	// Your teardown logic here
}

func TestExample(t *testing.T) {
	fmt.Println("Executing test")
	// Your test logic here
}

func TestMain(m *testing.M) {
	setup()
	exitCode := m.Run()
	teardown()
	os.Exit(exitCode)
}

func TestDefaultTimeForDB(t *testing.T) {

	var now = time.Now().UTC()

	var res = defaultTimeForDB()

	if res.Sub(now).Seconds() != 0 {
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
