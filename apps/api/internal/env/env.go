package env

import (
	"log"
	"os"
)

func Must(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("missing required %s\n", key)
	}

	return val
}

func Fallback(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return fallback
}
