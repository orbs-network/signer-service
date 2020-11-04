// Copyright 2019 the orbs-network-go authors
// This file is part of the orbs-network-go library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package config

import (
	"encoding/json"
	"github.com/orbs-network/crypto-lib-go/crypto/encoding"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

type ArrayFlags []string

func (i *ArrayFlags) String() string {
	return "my string representation"
}

func (i *ArrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func GetNodeConfigFromFiles(configFiles ArrayFlags, httpAddress string) (SignerServiceConfig, error) {
	cfg := make(map[string]string)

	if len(configFiles) != 0 {
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
