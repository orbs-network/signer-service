// Copyright 2019 the orbs-network-go authors
// This file is part of the orbs-network-go library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/orbs-network/signer-service/bootstrap/signer"
	"github.com/orbs-network/signer-service/config"
	"github.com/orbs-network/crypto-lib-go/crypto/ethereum/digest"
	crypto "github.com/orbs-network/crypto-lib-go/crypto/signer"
	"github.com/orbs-network/crypto-lib-go/test/crypto/ethereum/keys"
	"github.com/orbs-network/signer-service/test/with"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/scribe/log"
	"github.com/stretchr/testify/require"
	"math/rand"
	"net/http"
	"testing"
)

func randomAddressAndEndpoints() (address, httpEndpoint, httpSignEndpoint, httpSighEthEndpoint string) {
	address = fmt.Sprintf("localhost:%d", 8000+rand.Intn(5000))
	httpEndpoint = "http://"+address
	httpSignEndpoint = httpEndpoint+"/sign"
	httpSighEthEndpoint = httpEndpoint+"/eth-sign"

	return
}

func TestSignerServer(t *testing.T) {
	with.Context(func(ctx context.Context) {
		pk := keys.EcdsaSecp256K1KeyPairForTests(0).PrivateKey()
		nodeAddress := keys.EcdsaSecp256K1KeyPairForTests(0).NodeAddress()

		testLogger := log.GetLogger()

		address, httpEndpoint, _, _ := randomAddressAndEndpoints()
		server, err := signer.StartSignerServer(config.NewSignerServerConfig(address, nodeAddress, pk), testLogger)
		require.NoError(t, err)
		defer server.GracefulShutdown(ctx)

		c := crypto.NewSignerClient(httpEndpoint)

		payload := []byte("payload")

		signed, err := digest.SignAsNode(pk, payload)
		require.NoError(t, err)

		clientSignature, err := c.Sign(ctx, payload)
		require.NoError(t, err)

		require.EqualValues(t, signed, clientSignature)
	})
}

func TestSignerServerSignatureWithBadRequest(t *testing.T) {
	with.Context(func(ctx context.Context) {
		pk := keys.EcdsaSecp256K1KeyPairForTests(0).PrivateKey()
		nodeAddress := keys.EcdsaSecp256K1KeyPairForTests(0).NodeAddress()

		testLogger := log.GetLogger()

		address, _, httpSignEndpoint, _ := randomAddressAndEndpoints()
		server, err := signer.StartSignerServer(config.NewSignerServerConfig(address, nodeAddress, pk), testLogger)
		require.NoError(t, err)
		defer server.GracefulShutdown(ctx)

		resp, err := http.Post(httpSignEndpoint, "application/membuffers", bytes.NewBuffer(nil))
		require.NoError(t, err)

		require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestSignerServerEthSignatureWithBadRequest(t *testing.T) {
	with.Context(func(ctx context.Context) {
		pk := keys.EcdsaSecp256K1KeyPairForTests(0).PrivateKey()
		nodeAddress := keys.EcdsaSecp256K1KeyPairForTests(0).NodeAddress()

		testLogger := log.GetLogger()

		address, _, _, httpSignEthEndpoint := randomAddressAndEndpoints()
		server, err := signer.StartSignerServer(config.NewSignerServerConfig(address, nodeAddress, pk), testLogger)
		require.NoError(t, err)
		defer server.GracefulShutdown(ctx)

		resp, err := http.Post(httpSignEthEndpoint, "application/membuffers", bytes.NewBuffer(nil))
		require.NoError(t, err)

		require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)
	})
}


func TestSignerServerWithWrongConfiguration(t *testing.T) {
	with.Context(func(ctx context.Context) {
		pk := keys.EcdsaSecp256K1KeyPairForTests(0).PrivateKey()
		nodeAddress := primitives.NodeAddress("hello")

		testLogger := log.GetLogger()

		address, _, _, _ := randomAddressAndEndpoints()
		_, err := signer.StartSignerServer(config.NewSignerServerConfig(address, nodeAddress, pk), testLogger)
		require.EqualError(t, err, "node address a328846cd5b4979d68a8c58a9bdfeee657b34de7 derived from secret key does not match provided node address 68656c6c6f")
	})
}

func TestHealthcheck(t *testing.T) {
	with.Context(func(ctx context.Context) {
		pk := keys.EcdsaSecp256K1KeyPairForTests(0).PrivateKey()
		nodeAddress := keys.EcdsaSecp256K1KeyPairForTests(0).NodeAddress()

		testOutput := log.NewTestOutput(t, log.NewHumanReadableFormatter())
		testLogger := log.GetLogger().WithOutput(testOutput)

		address, httpEndpoint, _, _ := randomAddressAndEndpoints()
		server, err := signer.StartSignerServer(config.NewSignerServerConfig(address, nodeAddress, pk), testLogger)
		require.NoError(t, err)
		defer server.GracefulShutdown(ctx)

		response, err := http.Get(httpEndpoint)
		require.NoError(t, err)

		require.EqualValues(t, http.StatusOK, response.StatusCode)

		server.GracefulShutdown(ctx)
	})
}