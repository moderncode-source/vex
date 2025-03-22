// SPDX-FileCopyrightText: 2025 The Vex Authors.
//
// SPDX-License-Identifier: Apache-2.0 OR MIT
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. You may not use this file except in compliance with the
// terms of those licenses.

// Vex service.

package vex

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/net/netutil"
)

// Vex major, minor, and patch version numbers.
const (
	VersionMajor = 0
	VersionMinor = 0
	VersionPatch = 1
)

const (
	// Block clients from keeping connections open forever by setting a
	// deadline for reading request headers (effectively, connection's read
	// deadline, see [http.Server]).
	serverReadHeaderTimeout = 10 * time.Second

	// Maximum concurrent TCP connections that a server can accept.
	// We calculate it as: anticipated Request Rate * Request Duration.
	// Coupled together with [serverReadHeaderTimeout] to avoid waiting
	// indefinitely for clients that never close connections.
	serverMaxConnections = 50
)

// ErrNilServer is returned by [Service.Start] and [Service.Stop] if
// service's server is nil.
var ErrNilServer = errors.New("service's server must not be nil")

// Service defines parameters and provides functionality to run a Vex service.
// Use [New] to create a new valid service instance.
type Service struct {
	server *http.Server
	logger *zerolog.Logger
}

// New allocates and returns a new [Service] with [http.Server] that will
// listen on TCP network address addr and handle requests on incoming
// connections using [ServiceMux] handler. This is the recommended and default
// way to create a Vex service.
//
// To choose your own handler or fall back to [http.DefaultServeMux],
// use [NewWithHandler].
func New(addr string, logger *zerolog.Logger) *Service {
	return &Service{
		server: &http.Server{
			ReadHeaderTimeout: serverReadHeaderTimeout,
			Addr:              addr,
			Handler:           ServiceMux,
		},
		logger: logger,
	}
}

// NewWithHandler allocates and returns a new [Service] with [http.Server]
// that will listen on TCP network address addr and handle requests on
// incoming connections by calling [http.Server.Serve] with handler.
//
// If handler is nil, [http.DefaultServeMux] will be used.
//
// See also: [New].
func NewWithHandler(addr string, handler http.Handler, logger *zerolog.Logger) *Service {
	return &Service{
		server: &http.Server{
			ReadHeaderTimeout: serverReadHeaderTimeout,
			Addr:              addr,
			Handler:           handler,
		},
		logger: logger,
	}
}

// Start begins listening to and serving incoming requests to the service
// on the configured network address. Call [Service.Stop] to stop serving.
func (svc *Service) Start() error {
	if svc.server == nil {
		return ErrNilServer
	}

	l, err := net.Listen("tcp", svc.server.Addr)
	if err != nil {
		return fmt.Errorf("failed to start service: %v", err)
	}

	// Limit the number of concurrent connections to the service.
	ln := netutil.LimitListener(l, serverMaxConnections)

	err = svc.server.Serve(ln)
	if err == nil || err == http.ErrServerClosed {
		return nil
	}
	return fmt.Errorf("failed to serve service: %v", err)
}

// Stop gracefully shuts down the service. See [http.Server.Shutdown].
func (svc *Service) Stop(ctx context.Context) error {
	if svc.server == nil {
		return ErrNilServer
	}

	err := svc.server.Shutdown(ctx)
	if err == nil || err == http.ErrServerClosed {
		return nil
	}
	return fmt.Errorf("failed to stop service: %v", err)
}
