package util

import (
	"fmt"
	"os"
)

func GetEnvOrDie(key string) string {
	value := os.Getenv(key)
	if value == "" {
		fmt.Printf("missing req'd env var %q", key)
		os.Exit(1)
	}
	return value
}
