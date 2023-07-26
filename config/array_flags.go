// Copyright 2019 the orbs-network-go authors
// This file is part of the orbs-network-go library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/orbs-network/scribe/log"

	"github.com/orbs-network/crypto-lib-go/crypto/encoding"
	"github.com/pkg/errors"
)

type ArrayFlags []string

func (i *ArrayFlags) String() string {
	return "my string representation"
}

func (i *ArrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// Parse required node configuration from environment variables or config file.
// Environment variables override config file.
// If no configuration is found (meaning no environment variables set or config file passed), an error is returned
func GetNodeConfig(configFiles ArrayFlags, httpAddress string, logger log.Logger) (SignerServiceConfig, error) {
	cfg := make(map[string]string)

	nodeAddressEnvVar := os.Getenv("NODE_ADDRESS")
	privateKeyEnvVar := os.Getenv("NODE_PRIVATE_KEY")

	if nodeAddressEnvVar != "" && privateKeyEnvVar != "" {
		cfg["node-address"] = nodeAddressEnvVar
		cfg["node-private-key"] = privateKeyEnvVar
		logger.Info("Using environment variables for node configuration")

	} else if len(configFiles) != 0 {
		for _, configFile := range configFiles {
			if _, err := os.Stat(configFile); os.IsNotExist(err) {
				return nil, errors.Errorf("could not open config file: %s", err)
			}

			contents, err := ioutil.ReadFile(configFile)
			if err != nil {
				return nil, err
			}

			data := make(map[string]string)

			json.Unmarshal(contents, &data)
			for k, v := range data {
				cfg[k] = v
			}

			if err != nil {
				return nil, err
			}
		}
		logger.Info("Using config.json for node configuration")
	} else {
		return nil, errors.New("No configuration file or environment variables found for node config")
	}

	nodeAddress, err := encoding.DecodeHex(cfg["node-address"])
	if err != nil {
		return nil, errors.Wrapf(err, "could not decode node address")
	}

	nodePrivateKey, err := encoding.DecodeHex(cfg["node-private-key"])
	if err != nil {
		return nil, errors.Wrapf(err, "could not decode node private key")
	}

	result := NewSignerServerConfig(httpAddress, nodeAddress, nodePrivateKey)

	return result, nil
}
