package main

import (
	"fmt"
	"log"
	"net/http"
	app_config "test-devices-api/config"
	"test-devices-api/database"
	"test-devices-api/routes"

	"github.com/gorilla/mux"
)

const PORT = ":3001"

func main() {

	app_config.LoadConfig()
	//routes.SetRoutesAuth(os.Getenv(secrets.KEY_EXTERNAL_AUTH))
	database.Connect()

	route := mux.NewRouter()

	routes.AddApproutes(route)

	fmt.Println("Device App running port", PORT)

	log.Fatal(http.ListenAndServe(PORT, route))

}
