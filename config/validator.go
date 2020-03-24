package config

import (
	"bytes"
	"encoding/hex"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/signer-service/crypto/digest"
	"github.com/orbs-network/signer-service/crypto/signature"
	"github.com/pkg/errors"
)

func ValidateSigner(cfg SignerServiceConfig) error {
	if len(cfg.NodePrivateKey()) == 0 {
		return errors.New("node private key must not be empty")
	}
	return requireCorrectNodeAddressAndPrivateKey(cfg.NodeAddress(), cfg.NodePrivateKey())
}

func requireCorrectNodeAddressAndPrivateKey(address primitives.NodeAddress, key primitives.EcdsaSecp256K1PrivateKey) error {
	// FIXME make byte32
	msg := []byte{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}

	sign, err := signature.SignEcdsaSecp256K1(key, msg)
	if err != nil {
		return errors.Wrap(err, "could not create test sign")
	}

	recoveredPublicKey, err := signature.RecoverEcdsaSecp256K1(msg, sign)
	if err != nil {
		return errors.Wrap(err, "could not recover public key from test sign")
	}

	recoveredNodeAddress := digest.CalcNodeAddressFromPublicKey(recoveredPublicKey)
	if bytes.Compare(address, recoveredNodeAddress) != 0 {
		return errors.Errorf("node address %s derived from secret key does not match provided node address %s",
			hex.EncodeToString(recoveredNodeAddress), hex.EncodeToString(address))
	}
	return nil
}