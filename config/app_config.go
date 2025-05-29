package app_config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const key_external_auth = "EXTERNAL_AUTH"
const key_mongo_url = "MONGO_URL"
const key_db_name = "DB_NAME"

const DEFAULT_KEY_FOR_CONFIG = "APP_CONFIG"

var MONGO_URL = ""
var EXTERNAL_AUTH = ""
var DB_NAME = ""

var varLists = []string{key_external_auth, key_mongo_url, key_db_name}

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
		SetVars()
		return
	}

	currentEnvironment, ok := os.LookupEnv(DEFAULT_KEY_FOR_CONFIG)
	envPath := "../"

	fmt.Println("enviroment", currentEnvironment, ok)

	if len(currentEnvironment) == 0 || currentEnvironment == "PROD" {
		currentEnvironment = "PROD"
		envPath = "./"
	}

	err := godotenv.Load(envPath + currentEnvironment + ".env")

	if err != nil {
		panic(err)
	}

	for _, e := range varLists {

		value := os.Getenv(e)

		if len(value) == 0 {
			log.Fatal("missing " + e + " on .env or script")
		}

	}

	SetVars()

}

func SetVars() {

	MONGO_URL = os.Getenv(key_mongo_url)
	EXTERNAL_AUTH = os.Getenv(key_external_auth)
	DB_NAME = os.Getenv(key_db_name)

}
