package domain_devices

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const TagName string = "name"
const TagBrand string = "brand"
const TagState string = "state"
const TagModificationTime string = "modification_time"
const TagCreationTime string = "creation_time"
const DefaultInitState = DeviceStatusAvalaible

const MaxLenID = 50
const MaxLenBrand = 30
const MaxLenName = 45

type Devices struct {
	Name             string             `json:"name" bson:"name"`
	Brand            string             `json:"brand" bson:"brand"`
	State            DeviceState        `json:"state" bson:"state"`
	ModificationTime time.Time          `json:"modification_time" bson:"modification_time"`
	CreationTime     time.Time          `json:"creation_time" bson:"creation_time"`
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

type DeviceState string

const (
	DeviceStatusAvalaible DeviceState = "avalaible"
	DeviceStatusInUse     DeviceState = "in-use"
	DeviceStatusInactive  DeviceState = "inactive"
)

var DeviceStatesList = []string{string(DeviceStatusAvalaible), string(DeviceStatusInUse), string(DeviceStatusInactive)}
