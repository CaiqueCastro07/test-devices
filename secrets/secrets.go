package secrets

import (
	"log"
	"os"
	"test-devices-api/routes"

	"github.com/joho/godotenv"
)

const KEY_EXTERNAL_AUTH = "EXTERNAL_AUTH"
const KEY_MONGO_URL = "MONGO_URL"

var varLists = []string{KEY_EXTERNAL_AUTH, KEY_MONGO_URL}

func Load() {

	missing := false

	for _, e := range varLists {

		value := os.Getenv(e)

		if len(value) == 0 {
			missing = true
			break
		}

	}

	if !missing {
		return
	}

	erro := godotenv.Load()

	if erro != nil {
		log.Fatal("Variaveis de ambiente nao carregadas no runtime")
	}

	for _, e := range varLists {

		value := os.Getenv(e)

		if len(value) == 0 {
			log.Fatal("missing " + e + " on .env or script")
		}

	}

	routes.SetRoutesAuth(os.Getenv(KEY_EXTERNAL_AUTH))

}
