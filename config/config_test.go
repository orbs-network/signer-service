package config

import (
	"github.com/orbs-network/crypto-lib-go/crypto/encoding"
	"github.com/stretchr/testify/require"
	"testing"
)

var NODE_ADDRESS, _ = encoding.DecodeHex("6e2cb55e4cbe97bf5b1e731d51cc2c285d83cbf9")
var NODE_PRIVATE_KEY, _ = encoding.DecodeHex("426308c4d11a6348a62b4fdfb30e2cad70ab039174e2e8ea707895e4c644c4ec")
const HTTP_ADDRESS = "localhost:7777"

func Test_GetNodeConfigFromFiles(t *testing.T) {
	cfg, err := GetNodeConfigFromFiles([]string{"./test/keys.json"}, HTTP_ADDRESS)
	require.NoError(t, err)

	require.EqualValues(t, cfg.NodeAddress(), NODE_ADDRESS)
	require.EqualValues(t, cfg.NodePrivateKey(), NODE_PRIVATE_KEY)
	require.EqualValues(t, cfg.HttpAddress(), HTTP_ADDRESS)
}