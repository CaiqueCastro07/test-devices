package app_config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const key_external_auth = "EXTERNAL_AUTH"
const key_mongo_url = "MONGO_URL"
const key_db_name = "DB_NAME"
const key_port = "PORT"

var ENVIRONMENT = ""

const DEFAULT_KEY_FOR_CONFIG = "APP_CONFIG"

var MONGO_URL = ""
var EXTERNAL_AUTH = ""
var DB_NAME = ""
var PORT = ""

var varLists = []string{key_external_auth, key_mongo_url, key_db_name, key_port}

func LoadConfig() {

	missing := false

	for _, e := range varLists {

		value := os.Getenv(e)

		if len(value) == 0 {
			missing = true
			break
		}

	}

	if !missing {
		setVars()
		return
	}

	currentEnvironment, _ := os.LookupEnv(DEFAULT_KEY_FOR_CONFIG)

	if len(currentEnvironment) == 0 {
		currentEnvironment = "PROD"
	}

	ENVIRONMENT = currentEnvironment

	possiblePaths := []string{"", "./", "../", "../../", "../../../"}

	found := false

	for _, e := range possiblePaths {

		err := godotenv.Load(e + currentEnvironment + ".env")

		if err != nil {
			continue
		}

		found = true
		break

	}

	if !found {
		log.Fatal("could not find proper .env file")
	}

	for _, e := range varLists {

		value := os.Getenv(e)

		if len(value) == 0 {
			log.Fatal("missing " + e + " on .env or script")
		}

	}

	setVars()

}

func setVars() {

	MONGO_URL = os.Getenv(key_mongo_url)
	EXTERNAL_AUTH = os.Getenv(key_external_auth)
	DB_NAME = os.Getenv(key_db_name)
	PORT = os.Getenv(key_port)

}
