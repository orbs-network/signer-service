// Copyright 2019 the orbs-network-go authors
// This file is part of the orbs-network-go library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/orbs-network/scribe/log"
	"github.com/orbs-network/signer-service/bootstrap/signer"
	"github.com/orbs-network/signer-service/config"
)

func getLogger() log.Logger {
	return log.GetLogger().WithOutput(log.NewFormattingOutput(os.Stdout, log.NewJsonFormatter()))
}

func main() {
	httpAddressFlag := flag.String("listen", ":7777", "ip address and port for http server")
	version := flag.Bool("version", false, "returns information about version")

	var configFiles config.ArrayFlags
	flag.Var(&configFiles, "config", "path/to/config.json")

	flag.Parse()

	httpAddress := config.GetHttpAddress(*httpAddressFlag)

	if *version {
		fmt.Println(config.GetVersion())
		return
	}

	logger := getLogger()
	cfg, err := config.GetNodeConfig(configFiles, httpAddress, logger)
	if err != nil {
		logger.Error("failed to parse configuration", log.Error(err))
		os.Exit(1)
	}

	server, err := signer.StartSignerServer(cfg, logger)
	if err != nil {
		logger.Error("failed to start signer service", log.Error(err))
		os.Exit(1)
	}

	server.WaitUntilShutdown(context.Background())
}
