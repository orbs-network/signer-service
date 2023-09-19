package config

import (
	"os"
)

// Returns the `HTTP_ADDRESS` environment variable if it is set, otherwise the provided flag value
func GetHttpAddress(httpAddressFlag string) string {
	var httpAddress string

	if envVar := os.Getenv("HTTP_ADDRESS"); envVar != "" {
		httpAddress = envVar
	} else {
		httpAddress = httpAddressFlag
	}

	return httpAddress
}
