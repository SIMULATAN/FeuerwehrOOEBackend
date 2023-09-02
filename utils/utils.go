package utils

import (
	"log"
	"os"
)

func LoadEnvVariable(key string, reference *string) string {
	log.Println("Loading env variable", key)
	if *reference == "" {
		*reference = os.Getenv(key)
	}
	log.Println("Loaded env variable", key, "with value", *reference)
	return *reference
}
