package config

import (
	"os"
	"testing"
)

func TestGetHttpAddress(t *testing.T) {
	// Default value for this flag
	httpAddressFlag := ":7777"

	// Environment variable should override the httpAddressFlag flag
	envVarName := "HTTP_ADDRESS"
	envVarValue := "localhost:8000"
	os.Setenv(envVarName, envVarValue)
	expected := envVarValue
	if result := GetHttpAddress(httpAddressFlag); result != expected {
		t.Errorf("Expected %v but got %v", expected, result)
	}
	os.Unsetenv(envVarName)

	// When HTTP_ADDRESS environment variable is not set and a flag is provided, it should return the provided flag value
	flagValues := []string{":7777", "localhost:80"}

	for _, flagValue := range flagValues {
		if result := GetHttpAddress(flagValue); result != flagValue {
			t.Errorf("Expected %v but got %v", flagValue, result)
		}
	}
}
