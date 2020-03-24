package config

import "github.com/orbs-network/orbs-spec/types/go/primitives"

type SignerConfig interface {
	NodePrivateKey() primitives.EcdsaSecp256K1PrivateKey
	SignerEndpoint() string
}

type SignerServiceConfig interface {
	NodeAddress() primitives.NodeAddress
	NodePrivateKey() primitives.EcdsaSecp256K1PrivateKey
	HttpAddress() string
}

type HttpServerConfig interface {
	HttpAddress() string
	Profiling() bool
}

type signerServerConfig struct {
	privateKey  primitives.EcdsaSecp256K1PrivateKey
	nodeAddress primitives.NodeAddress
	httpAddress string
}

func (s *signerServerConfig) NodePrivateKey() primitives.EcdsaSecp256K1PrivateKey {
	return s.privateKey
}

func (s *signerServerConfig) HttpAddress() string {
	return s.httpAddress
}

func (s *signerServerConfig) NodeAddress() primitives.NodeAddress {
	return s.nodeAddress
}

func (s *signerServerConfig) Profiling() bool {
	return false;
}

func NewSignerServerConfig(httpAddress string, nodeAddress primitives.NodeAddress, privateKey primitives.EcdsaSecp256K1PrivateKey) SignerServiceConfig {
	return &signerServerConfig{
		httpAddress: httpAddress,
		nodeAddress: nodeAddress,
		privateKey:  privateKey,
	}
}

