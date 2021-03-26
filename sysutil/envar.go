// +build admin

package sysutil

import (
	"log"
	"os"
)

func GetRequiredEnvString(key string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		log.Fatalf("Error: Missing environment variable %v", key)
	}
	return val
}

func GetEnvStringOrDefault(key, def string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return def
}
