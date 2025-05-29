package routes

import (
	"encoding/json"
	"net/http"
	"strings"
	domain_devices "test-devices-api/domain"
	usecase_devices "test-devices-api/useCase/devices"
	"time"

	"github.com/gorilla/mux"
)

const (
	ErrorInvalidBodyCreateRoute RoutesErrors = `invalid body sent, must be JSON = {"brand": "string","name":"string" }`
	ErrorInvalidIDParam         RoutesErrors = `invalid id param at the end of the URL path`
)

var ErrorInvalidBodyUpdateRoute string = `invalid body sent, must be JSON = {"brand": "string","name":"string", "state":"` + strings.Join(domain_devices.DeviceStatesList[:], " or ") + `"}`

type PingStatusResponse struct {
	Time   string `json:"time"`
	Status string `json:"status"`
}

func Ping(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(200)
	json.NewEncoder(res).Encode(&PingStatusResponse{
		Time:   time.Now().UTC().Format(time.RFC3339),
		Status: "UP",
	})

}

type CreateResponse struct {
	RecordId string `json:"record_id"`
}

type DefaultErrorResponse struct {
	Error string `json:"error"`
}

func CreateDeviceReq(res http.ResponseWriter, req *http.Request) {

	req.Body = http.MaxBytesReader(res, req.Body, maxBodyBytesSize)

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()

	defer req.Body.Close()

	var dados domain_devices.Devices

	errDecodeBody := dec.Decode(&dados)

	if errDecodeBody != nil {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(&DefaultErrorResponse{
			Error: string(ErrorInvalidBodyCreateRoute),
		})
		return
	}

	createdId, errorCreation := usecase_devices.Create(dados)

	if errorCreation != nil {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(&DefaultErrorResponse{
			Error: errorCreation.Error(),
		})
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(&CreateResponse{
		RecordId: createdId,
	})

}

type UpdateErrorResponse struct {
	Error string `json:"error"`
}

func UpdateDeviceReq(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]

	if len(id) == 0 || len(id) > domain_devices.MaxLenID {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(&DefaultErrorResponse{
			Error: string(ErrorInvalidIDParam),
		})
		return
	}

	req.Body = http.MaxBytesReader(res, req.Body, maxBodyBytesSize)

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()

	defer req.Body.Close()

	var dados domain_devices.Devices

	errDecodeBody := dec.Decode(&dados)

	if errDecodeBody != nil {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(&DefaultErrorResponse{
			Error: string(ErrorInvalidBodyUpdateRoute),
		})
		return
	}

	errorUpdate := usecase_devices.Update(id, dados)

	if errorUpdate != nil {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(&DefaultErrorResponse{
			Error: errorUpdate.Error(),
		})
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

}

func GetDeviceByIDReq(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]

	if len(id) == 0 || len(id) > domain_devices.MaxLenID {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(&DefaultErrorResponse{
			Error: string(ErrorInvalidIDParam),
		})
		return
	}

	device, errorGetDevice := usecase_devices.GetDeviceById(id)

	if errorGetDevice != nil {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(&DefaultErrorResponse{
			Error: errorGetDevice.Error(),
		})
		return
	}

	if len(device.Name) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(device)

}

func GetAllDevicesReq(res http.ResponseWriter, req *http.Request) {
	// paging
	brandFilter := req.URL.Query().Get(domain_devices.TagBrand)
	nameFilter := req.URL.Query().Get(domain_devices.TagName)
	stateFilter := req.URL.Query().Get(domain_devices.TagState)

	var filters = domain_devices.Devices{
		Brand: brandFilter,
		Name:  nameFilter,
		State: domain_devices.DeviceState(stateFilter),
	}

	devices, errorGetDevices := usecase_devices.GetAllDevicesWithFilters(filters)

	if errorGetDevices != nil {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(&DefaultErrorResponse{
			Error: errorGetDevices.Error(),
		})
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(devices)

}

func DeleteDeviceByIDReq(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]

	if len(id) == 0 || len(id) > domain_devices.MaxLenID {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(&DefaultErrorResponse{
			Error: string(ErrorInvalidIDParam),
		})
		return
	}

	errorDeleteDevice := usecase_devices.DeleteDeviceById(id)

	if errorDeleteDevice != nil {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(&DefaultErrorResponse{
			Error: errorDeleteDevice.Error(),
		})
		return
	}

	res.WriteHeader(http.StatusNoContent)

}
