// Copyright 2019 the orbs-network-go authors
// This file is part of the orbs-network-go library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package signer

import (
	"context"
	"github.com/orbs-network/govnr"
	"github.com/orbs-network/scribe/log"
	"net"
	"net/http"
)

type httpServer struct {
	govnr.ShutdownWaiter
	server *http.Server
	port   int
	logger log.Logger
	router *http.ServeMux
}

// TODO: unify with httpserver.HttpServer
func NewHttpServer(address string, logger log.Logger) (*httpServer, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	logger.Info("started http server", log.String("address", address))

	router := http.NewServeMux()

	s := &httpServer{
		server: &http.Server{
			Handler: router,
		},
		port:   listener.Addr().(*net.TCPAddr).Port,
		logger: logger,
		router: router,
	}

	go s.server.Serve(listener)

	return s, nil
}

func (s *httpServer) GracefulShutdown(shutdownContext context.Context) {
	if err := s.server.Shutdown(shutdownContext); err != nil {
		s.logger.Error("failed to stop http HttpServer gracefully", log.Error(err))
	}

}

func (s *httpServer) Port() int {
	return s.port
}

func (s *httpServer) Router() *http.ServeMux {
	return s.router
}

func (c *httpServer) WaitUntilShutdown(shutdownContext context.Context) {
	select {
	case <-shutdownContext.Done():
		if shutdownContext.Err() == context.DeadlineExceeded {
			panic("failed to shutdown")
		}
	}
}
