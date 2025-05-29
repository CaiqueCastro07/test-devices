package database

import (
	"context"
	"errors"
	"log"
	app_config "test-devices-api/config"
	domain_devices "test-devices-api/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Database

const devices_collection string = "devices"
const defaultIDTagName string = "_id"

type mongoOperators string

const (
	queryANDoperator mongoOperators = "$and"
	queryORoperator  mongoOperators = "$or"
	querySEToperator mongoOperators = "$set"
)

func defaultTimeForDB() time.Time {
	return time.Now().UTC()
}

type DBStringErrors string

const (
	ErrorDBNotConnected DBStringErrors = "the database is not connected"
	ErrorEmptyUpdateMap DBStringErrors = "empty update map"
	ErrorEmptySearchMap DBStringErrors = "empty search map"
)

func Connect() {
	// "www.mongodb.com/docs/drivers/go/current/"
	mongoURL := app_config.MONGO_URL

	if len(mongoURL) == 0 {
		log.Fatal("URL do mongoDB não setada corretamente no .env")
	}

	client, err := mongo.Connect(context.Background(), options.Client().
		ApplyURI(mongoURL))

	if err != nil {
		log.Fatal(err)
	}

	dbClient = client.Database(app_config.DB_NAME)

}

func InsertDevice(deviceData domain_devices.Devices) (string, error) {

	if dbClient == nil {
		return "", errors.New(string(ErrorDBNotConnected))
	}

	startTime := defaultTimeForDB()

	deviceData.CreationTime = startTime
	deviceData.ModificationTime = startTime

	deviceData.State = domain_devices.DefaultInitState

	res, dbError := dbClient.Collection(devices_collection).InsertOne(context.Background(), deviceData)

	if dbError != nil {
		return "", dbError
	}

	id, ok := res.InsertedID.(primitive.ObjectID)

	if !ok {
		return "", errors.New("the insertion id did not return correctly")
	}

	return id.Hex(), nil

}

func UpdateDeviceByID(deviceID string, updateMap domain_devices.Devices) error {

	if dbClient == nil {
		return errors.New(string(ErrorDBNotConnected))
	}

	updateSet := bson.M{}

	if len(updateMap.Name) > 0 {
		updateSet[domain_devices.TagName] = updateMap.Name
	}

	if len(updateMap.Brand) > 0 {
		updateSet[domain_devices.TagBrand] = updateMap.Brand
	}

	if len(updateMap.State) > 0 {
		updateSet[domain_devices.TagState] = updateMap.State
	}

	if len(updateSet) == 0 {
		return errors.New(string(ErrorEmptyUpdateMap))
	}

	objID, erroId := primitive.ObjectIDFromHex(deviceID)

	if erroId != nil {
		return erroId
	}

	updateSet[domain_devices.TagModificationTime] = defaultTimeForDB()

	update := bson.M{
		string(querySEToperator): updateSet,
	}

	dbError := dbClient.Collection(devices_collection).FindOneAndUpdate(context.Background(), bson.M{defaultIDTagName: objID}, update)

	if dbError.Err() != nil {
		return dbError.Err()
	}

	return nil

}

func GetDeviceByID(id string) (*domain_devices.Devices, error) {

	if dbClient == nil {
		return &domain_devices.Devices{}, errors.New(string(ErrorDBNotConnected))
	}

	objID, erroId := primitive.ObjectIDFromHex(id)

	if erroId != nil {
		return &domain_devices.Devices{}, erroId
	}

	var result domain_devices.Devices

	err := dbClient.Collection(devices_collection).FindOne(context.Background(), bson.D{{defaultIDTagName, objID}}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		return &domain_devices.Devices{}, nil
	}

	if err != nil {
		return &domain_devices.Devices{}, err
	}

	return &result, nil

}

func GetDeviceByStateID(id string) (domain_devices.DeviceState, error) {

	if dbClient == nil {
		return "", errors.New(string(ErrorDBNotConnected))
	}

	objID, erroId := primitive.ObjectIDFromHex(id)

	if erroId != nil {
		return "", erroId
	}

	fieldsToBring := options.FindOne().SetProjection(bson.D{{domain_devices.TagState, 1}})

	var result domain_devices.Devices

	err := dbClient.Collection(devices_collection).FindOne(context.Background(), bson.D{{defaultIDTagName, objID}}, fieldsToBring).Decode(&result)

	if err == mongo.ErrNoDocuments {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return result.State, nil

}

func GetDeviceByField(searchMap domain_devices.Devices) (*domain_devices.Devices, error) {

	if dbClient == nil {
		return &domain_devices.Devices{}, errors.New(string(ErrorDBNotConnected))
	}

	filterMap := []interface{}{}

	if len(searchMap.State) > 0 {
		filterMap = append(filterMap, bson.D{{domain_devices.TagState, searchMap.State}})
	}

	if len(searchMap.Brand) > 0 {
		filterMap = append(filterMap, bson.D{{domain_devices.TagBrand, searchMap.Brand}})
	}

	if len(searchMap.Name) > 0 {
		filterMap = append(filterMap, bson.D{{domain_devices.TagName, searchMap.Name}})
	}

	if len(filterMap) == 0 {
		return &domain_devices.Devices{}, errors.New(string(ErrorEmptySearchMap))
	}

	filter := bson.D{
		{string(queryANDoperator), filterMap},
	}

	var device domain_devices.Devices

	err := dbClient.Collection(devices_collection).FindOne(context.Background(), filter).Decode(&device)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			return &domain_devices.Devices{}, nil
		}

		return &domain_devices.Devices{}, err
	}

	return &device, nil

}

func GetAllDevicesByField(searchMap domain_devices.Devices) (*[]domain_devices.Devices, error) {

	if dbClient == nil {
		return &[]domain_devices.Devices{}, errors.New(string(ErrorDBNotConnected))
	}

	filterMap := []interface{}{}

	if len(searchMap.State) > 0 {
		filterMap = append(filterMap, bson.D{{domain_devices.TagState, searchMap.State}})
	}

	if len(searchMap.Brand) > 0 {
		filterMap = append(filterMap, bson.D{{domain_devices.TagBrand, searchMap.Brand}})
	}

	if len(searchMap.Name) > 0 {
		filterMap = append(filterMap, bson.D{{domain_devices.TagName, searchMap.Name}})
	}

	var filter = bson.D{}

	if len(filterMap) > 0 {
		filter = bson.D{
			{string(queryANDoperator), filterMap},
		}
	}

	res, err := dbClient.Collection(devices_collection).Find(context.Background(), filter)

	if err != nil {
		return &[]domain_devices.Devices{}, err
	}

	defer res.Close(context.Background())

	var list = []domain_devices.Devices{}

	erro := res.All(context.TODO(), &list)

	if erro != nil {
		return &[]domain_devices.Devices{}, erro
	}

	return &list, nil

}

func DeleteDeviceByID(id string) error {

	if dbClient == nil {
		return errors.New(string(ErrorDBNotConnected))
	}

	objID, erroId := primitive.ObjectIDFromHex(id)

	if erroId != nil {
		return erroId
	}

	var result domain_devices.Devices

	err := dbClient.Collection(devices_collection).FindOneAndDelete(context.Background(), bson.D{{defaultIDTagName, objID}}).Decode(&result)

	if err != nil {
		return err
	}

	return nil

}

func DeleteAllDevices() error {

	if dbClient == nil {
		return errors.New(string(ErrorDBNotConnected))
	}

	_, err := dbClient.Collection(devices_collection).DeleteMany(context.Background(), bson.D{{}})

	if err != nil {
		return err
	}

	return nil

}
