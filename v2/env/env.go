package env

import "os"

func GetEnv(name string, fallback string) string {
	if val := os.Getenv(name); len(val) > 0 {
		return val
	} else {
		return fallback
	}
}
