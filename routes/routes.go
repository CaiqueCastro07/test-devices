package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func setStaticFolder(route *mux.Router) {
	fs := http.FileServer(http.Dir("./public/"))
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))
}

type RoutesPaths string
type RoutesErrors string

const (
	DevicesPath string = "/devices"
	StatusPath  string = "/status"
)

var authorizedPaths = map[string]bool{
	StatusPath: true,
}

var auth = ""

func SetRoutesAuth(token string) {
	if len(token) > 0 {
		auth = token
	}
}

func isValidAuth(userAuth string) bool {
	return userAuth == auth
}

const headerForAuth = "x-api-key"
const maxBodyBytesSize int64 = 300000

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		if authorizedPaths[req.RequestURI] {
			next.ServeHTTP(res, req)
			return
		}

		authorized := isValidAuth(req.Header.Get(headerForAuth))

		if !authorized {
			res.WriteHeader(401)
			return
		}

		next.ServeHTTP(res, req)

	})
}

func AddApproutes(route *mux.Router) {

	setStaticFolder(route)

	route.Use(authMiddleware)

	route.HandleFunc(StatusPath, Ping).Methods("GET")

	route.HandleFunc(DevicesPath, GetAllDevicesReq).Methods("GET")
	route.HandleFunc(DevicesPath+"/{id}", GetDeviceByIDReq).Methods("GET")
	route.HandleFunc(DevicesPath, CreateDeviceReq).Methods("POST")
	route.HandleFunc(DevicesPath+"/{id}", UpdateDeviceReq).Methods("PUT")
	route.HandleFunc(DevicesPath+"/{id}", DeleteDeviceByIDReq).Methods("DELETE")

}
