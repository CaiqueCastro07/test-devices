package main

import (
	"fmt"
	"log"
	"net/http"
	"test-devices-api/database"
	"test-devices-api/routes"
	"test-devices-api/secrets"

	"github.com/gorilla/mux"
)

const PORT = ":3001"

func main() {

	secrets.Load()

	database.Connect()

	route := mux.NewRouter()

	routes.AddApproutes(route)

	fmt.Println("Device App running port", PORT)

	log.Fatal(http.ListenAndServe(PORT, route))

}
