package config

import (
	"os"
	"testing"

	"github.com/orbs-network/crypto-lib-go/crypto/encoding"
	"github.com/orbs-network/scribe/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const HTTP_ADDRESS = "localhost:7777"

// Ensure configuration can be set via a config file
func TestGetNodeConfig(t *testing.T) {
	var NODE_ADDRESS, _ = encoding.DecodeHex("6e2cb55e4cbe97bf5b1e731d51cc2c285d83cbf9")
	var NODE_PRIVATE_KEY, _ = encoding.DecodeHex("426308c4d11a6348a62b4fdfb30e2cad70ab039174e2e8ea707895e4c644c4ec")

	testLogger := log.GetLogger()
	cfg, err := GetNodeConfig([]string{"./test/keys.json"}, HTTP_ADDRESS, testLogger)
	require.NoError(t, err)

	require.EqualValues(t, cfg.NodeAddress(), NODE_ADDRESS)
	require.EqualValues(t, cfg.NodePrivateKey(), NODE_PRIVATE_KEY)
	require.EqualValues(t, cfg.HttpAddress(), HTTP_ADDRESS)
}

// Ensure configuration can be set via environment variables
func TestGetConfigFromEnvVars(t *testing.T) {
	var configFiles ArrayFlags

	addressEnvVarName := "NODE_ADDRESS"
	addressEnvVarValue := "555cb55e4cbe97bf5b1e731d51cc2c285d83cbf9"
	os.Setenv(addressEnvVarName, addressEnvVarValue)
	defer os.Unsetenv(addressEnvVarName)

	privateKeyEnvVarName := "NODE_PRIVATE_KEY"
	privateKeyEnvVarNameValue := "555308c4d11a6348a62b4fdfb30e2cad70ab039174e2e8ea707895e4c644c4ec"
	os.Setenv(privateKeyEnvVarName, privateKeyEnvVarNameValue)
	defer os.Unsetenv(privateKeyEnvVarName)

	testLogger := log.GetLogger()
	cfg, err := GetNodeConfig(configFiles, HTTP_ADDRESS, testLogger)

	require.NoError(t, err)

	NODE_ADDRESS, _ := encoding.DecodeHex(addressEnvVarValue)
	NODE_PRIVATE_KEY, _ := encoding.DecodeHex(privateKeyEnvVarNameValue)

	require.EqualValues(t, cfg.NodeAddress(), NODE_ADDRESS)
	require.EqualValues(t, cfg.NodePrivateKey(), NODE_PRIVATE_KEY)
	require.EqualValues(t, cfg.HttpAddress(), HTTP_ADDRESS)
}

// Ensure environment variables take precedence over config file
func TestEnvVarsTakePrecedence(t *testing.T) {

	addressEnvVarName := "NODE_ADDRESS"
	addressEnvVarValue := "555cb55e4cbe97bf5b1e731d51cc2c285d83cbf9"
	os.Setenv(addressEnvVarName, addressEnvVarValue)
	defer os.Unsetenv(addressEnvVarName)

	privateKeyEnvVarName := "NODE_PRIVATE_KEY"
	privateKeyEnvVarNameValue := "555308c4d11a6348a62b4fdfb30e2cad70ab039174e2e8ea707895e4c644c4ec"
	os.Setenv(privateKeyEnvVarName, privateKeyEnvVarNameValue)
	defer os.Unsetenv(privateKeyEnvVarName)

	testLogger := log.GetLogger()
	cfg, err := GetNodeConfig([]string{"./test/keys.json"}, HTTP_ADDRESS, testLogger)

	require.NoError(t, err)

	NODE_ADDRESS, _ := encoding.DecodeHex(addressEnvVarValue)
	NODE_PRIVATE_KEY, _ := encoding.DecodeHex(privateKeyEnvVarNameValue)

	require.EqualValues(t, cfg.NodeAddress(), NODE_ADDRESS)
	require.EqualValues(t, cfg.NodePrivateKey(), NODE_PRIVATE_KEY)
}

// Ensure an error is returned when:
// 1. No configuration file or environment variables are set
// 2. Only one environment variable is set
func TestGetConfigNoFileOrEnvVars(t *testing.T) {
	errorMsg := "No configuration file or environment variables found for node config"
	testLogger := log.GetLogger()

	t.Run("noFileOrEnvVars", func(t *testing.T) {
		var configFiles ArrayFlags
		cfg, err := GetNodeConfig(configFiles, HTTP_ADDRESS, testLogger)

		assert.Nil(t, cfg)
		require.Errorf(t, err, errorMsg)
	})

	t.Run("onlyOneEnvVarSet", func(t *testing.T) {
		var configFiles ArrayFlags
		// Set only one env var
		addressEnvVarName := "NODE_ADDRESS"
		addressEnvVarValue := "444cb55e4cbe97bf5b1e731d51cc2c285d83cbf9"
		os.Setenv(addressEnvVarName, addressEnvVarValue)
		defer os.Unsetenv(addressEnvVarName)

		cfg, err := GetNodeConfig(configFiles, HTTP_ADDRESS, testLogger)

		assert.Nil(t, cfg)
		require.Errorf(t, err, "No configuration file or environment variables found for node config")
	})
}
