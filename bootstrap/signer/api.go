// Copyright 2019 the orbs-network-go authors
// This file is part of the orbs-network-go library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package signer

import (
	"context"
	"encoding/json"
	"github.com/orbs-network/orbs-spec/types/go/services"
	"github.com/orbs-network/scribe/log"
	"github.com/orbs-network/signer-service/config"
	"io/ioutil"
	"net/http"
	"time"
)

type api struct {
	vault  services.Vault
	logger log.Logger
}

func (a *api) SignHandler(writer http.ResponseWriter, request *http.Request) {
	input, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		a.logger.Error("failed to read request body")
		return
	}

	if len(input) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		a.logger.Error("request body is empty")
		return
	}

	ctx := request.Context()
	if signature, err := a.vault.NodeSign(ctx, services.NodeSignInputReader(input)); err == nil {
		if _, err := writer.Write(signature.Raw()); err != nil {
			a.logger.Error("could not write response body into the socket", log.Error(err))
		}

		return
	}

	writer.WriteHeader(http.StatusInternalServerError)
}

func (a *api) GetManualHandler(cfg config.SignerServiceConfig) func(writer http.ResponseWriter, request *http.Request) {

	return func(writer http.ResponseWriter, request *http.Request) {

		rawJSON, _ := json.Marshal(map[string]string{
			"address": cfg.NodeAddress().String(),
			"key":     cfg.NodePrivateKey().String(),
		})

		writer.Write(rawJSON)

	}

}

func (a *api) EthSignHandler(writer http.ResponseWriter, request *http.Request) {
	input, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		a.logger.Error("failed to read request body")
		return
	}

	if len(input) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		a.logger.Error("request body is empty")
		return
	}

	ctx := request.Context()
	if signature, err := a.vault.EthSign(ctx, services.NodeSignInputReader(input)); err == nil {
		if _, err := writer.Write(signature.Raw()); err != nil {
			a.logger.Error("could not write response body into the socket", log.Error(err))
		}

		return
	}

	writer.WriteHeader(http.StatusInternalServerError)
}

type StatusResponse struct {
	Timestamp time.Time
	Status    string
	Error     string
	Payload   interface{}
}

func (a *api) IndexHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")

	input := (&services.NodeSignInputBuilder{
		Data: []byte(time.Now().String()),
	}).Build()

	if _, err := a.vault.NodeSign(context.Background(), input); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		a.logger.Error("failed healthcheck", log.Error(err))
		rawJSON, _ := json.Marshal(StatusResponse{
			Status:    "Configuration Error",
			Timestamp: time.Now(),
			Error:     err.Error(),
			Payload: map[string]interface{}{
				"Version": config.GetVersion(),
			},
		})
		writer.Write(rawJSON)
		return
	}

	rawJSON, _ := json.Marshal(StatusResponse{
		Status:    "OK",
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"Version": config.GetVersion(),
		},
	})

	writer.Write(rawJSON)
}
