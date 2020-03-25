package config

import (
	"bytes"
	"encoding/hex"
	"github.com/orbs-network/crypto-lib-go/crypto/digest"
	"github.com/orbs-network/crypto-lib-go/crypto/signature"
	"github.com/pkg/errors"
)

func ValidateSigner(cfg SignerServiceConfig) error {
	if len(cfg.NodePrivateKey()) == 0 {
		return errors.New("node private key must not be empty")
	}
	address := cfg.NodeAddress()
	key := cfg.NodePrivateKey()

	msg := make([]byte, 32)
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