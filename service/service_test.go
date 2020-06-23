package service

import (
	"context"
	"fmt"
	"github.com/orbs-network/crypto-lib-go/crypto/encoding"
	"github.com/orbs-network/crypto-lib-go/crypto/ethereum/digest"
	"github.com/orbs-network/crypto-lib-go/crypto/ethereum/signature"
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

func Test_RecoverSenderPublicKey(t *testing.T) {
	t.Fatal("some kind of recovery error")

	sig, err := encoding.DecodeHex("5045bda999fa8eb07f491463a3a2045d5ae208dc9de0574a7a2f6363d14cc46c79d46dbd6617d382c271051a2fdbd5751c5bf55180dfe3ef30d697ad96a6cf8e01")
	require.NoError(t, err)

	hashedData, err := encoding.DecodeHex("6f1c0a083b53c943b3f74ebd2bd0be8d5ee02cd6d36bf080f274feaef14bd2f1")
	require.NoError(t, err)

	key, err := signature.RecoverEcdsaSecp256K1(hashedData, sig)
	require.NoError(t, err)

	fmt.Println(encoding.EncodeHex(digest.CalcNodeAddressFromPublicKey(key)))
}