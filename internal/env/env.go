package env

import (
	"log"
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)

	if !ok {
		log.Printf("Warning: Environment variable %s not set, using fallback: %s", key, fallback)
		return fallback
	}

	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)

	if !ok {
		log.Printf("Warning: Environment variable %s not set, using fallback: %d", key, fallback)
		return fallback
	}

	intVal, err := strconv.Atoi(val)

	if err != nil {
		log.Printf("Error: Invalid value for %s: %s, using fallback: %d", key, val, fallback)
		return fallback
	}

	return intVal
}
