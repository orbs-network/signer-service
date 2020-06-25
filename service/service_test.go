package service

import (
	"context"
	"github.com/orbs-network/crypto-lib-go/crypto/encoding"
	"github.com/orbs-network/crypto-lib-go/crypto/ethereum/digest"
	"github.com/orbs-network/crypto-lib-go/crypto/ethereum/signature"
	"github.com/orbs-network/crypto-lib-go/crypto/hash"
	"github.com/orbs-network/crypto-lib-go/test/crypto/ethereum/keys"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/services"
	"github.com/orbs-network/scribe/log"
	"github.com/stretchr/testify/require"
	"testing"
)

type signerServiceConfig struct {
}

func (s *signerServiceConfig) NodePrivateKey() primitives.EcdsaSecp256K1PrivateKey {
	return keys.EcdsaSecp256K1KeyPairForTests(0).PrivateKey()
}

func TestService_NodeSign(t *testing.T) {
	cfg := &signerServiceConfig{}
	pk := cfg.NodePrivateKey()

	testOutput := log.NewTestOutput(t, log.NewHumanReadableFormatter())
	testLogger := log.GetLogger().WithOutput(testOutput)

	service := NewService(cfg, testLogger)

	payload := []byte("payload")

	signed, err := digest.SignAsNode(pk, payload)
	require.NoError(t, err)

	clientSignature, err := service.NodeSign(context.TODO(), (&services.NodeSignInputBuilder{
		Data: payload,
	}).Build())
	require.NoError(t, err)

	require.EqualValues(t, signed, clientSignature.Signature())
}

// Contract values for JS
func Test_Payload(t *testing.T) {
	payload := []byte{1, 2, 3}

	raw := (&services.NodeSignInputBuilder{
		Data: payload,
	}).Build().Raw()

	require.EqualValues(t,[]byte{3, 0, 0, 0, 1, 2, 3}, raw)
}

func TestService_EthSign(t *testing.T) {
	cfg := &signerServiceConfig{}
	pk := cfg.NodePrivateKey()

	testOutput := log.NewTestOutput(t, log.NewHumanReadableFormatter())
	testLogger := log.GetLogger().WithOutput(testOutput)

	service := NewService(cfg, testLogger)

	payload, err := encoding.DecodeHex("f83f808b32303030303030303030308532313030309435353535353535353535353535353535353535359331303030303030303030303030303030303030820001")
	require.NoError(t, err)

	h := hash.CalcKeccak256(payload)
	signed, err := signature.SignEcdsaSecp256K1(pk, h)
	require.NoError(t, err)

	require.EqualValues(t, "50b69b24790fbdf91bd0272fef54f7490fb4f61cb07a91a3d61e6c115a6fe80b76df08f4f3a5763bc721423c89c074fec9af0ed86bf889973a85499c4691cbf201", signed.String())

	clientSignature, err := service.EthSign(context.TODO(), (&services.NodeSignInputBuilder{
		Data: payload,
	}).Build())
	require.NoError(t, err)

	require.EqualValues(t, signed, clientSignature.Signature())
}