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

// ErrNilServer is returned by [Service.Validate], [Service.Start] and
// [Service.Stop] if service's server is nil.
var ErrNilServer = errors.New("service's server must not be nil")

// ErrNilLogger is returned by [Service.Validate] if service's logger is nil.
var ErrNilLogger = errors.New("service's logger must not be nil")

// Service defines parameters and provides functionality to run a Vex service.
// Use [New] to create a new valid service instance.
type Service struct {
	server *http.Server
	logger *zerolog.Logger
}

// New allocates and returns a new [Service] with [http.Server] that will
// listen on TCP network address addr and handle requests on incoming
// connections using [ServiceMux] handler. This is the recommended and the
// default way to create a Vex service.
//
// Default service request handlers will be registered automatically. If you
// will need multiple Vex services, use [NewWithHandler] instead and pass a new
// handler to each to avoid conflict errors.
//
// To choose your own handler or fall back to [http.DefaultServeMux],
// use [NewWithHandler].
func New(addr string, logger *zerolog.Logger) (*Service, error) {
	svc := &Service{
		server: &http.Server{
			ReadHeaderTimeout: serverReadHeaderTimeout,
			Addr:              addr,
			Handler:           ServiceMux,
		},
		logger: logger,
	}

	if err := svc.RegisterDefaultHandlers(ServiceMux); err != nil {
		return nil, err
	}

	return svc, nil
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

// Validate checks whether the service is in a
// valid state and its parameters are valid.
func (svc *Service) Validate() error {
	if svc.server == nil {
		return ErrNilServer
	}

	if svc.logger == nil {
		return ErrNilLogger
	}

	if _, err := net.ResolveTCPAddr("tcp", svc.server.Addr); err != nil {
		return fmt.Errorf("failed to resolve server addr: %v", err)
	}

	return nil
}

// RegisterDefaultHandlers registers all default request handlers for the
// service on the given HTTP request multiplexer. If [http.ServeMux.HandleFunc]
// panics, its error is captured and returned.
func (svc *Service) RegisterDefaultHandlers(mux *http.ServeMux) (err error) {
	// Let the caller handle the error if [http.ServeMux.HandleFunc] panics.
	defer func() {
		if recoverErr, ok := recover().(error); ok {
			err = recoverErr
		}
	}()

	mux.HandleFunc(HealthEndpoint, svc.HealthHandler)

	// Request handlers' endpoints for the mux below start with "/v1/".
	//
	// We could instead create another mux with a handler wrapped in
	// [http.StripPrefix] to make endpoint patterns shorter, but, since there
	// is a small total number of endpoints, it is unnecessary.
	mux.HandleFunc(ReadyEndpoint, svc.ReadyHandler)

	// Submission queue GET/POST handlers.
	mux.HandleFunc(PostQueueEndpoint, svc.PostQueueHandler)
	mux.HandleFunc(GetQueueEndpoint, svc.GetQueueHandler)

	return err
}

// Start begins listening to and serving incoming requests to the service
// on the configured network address. Call [Service.Stop] to stop serving.
func (svc *Service) Start() error {
	if err := svc.Validate(); err != nil {
		return err
	}

	l, err := net.Listen("tcp", svc.server.Addr)
	if err != nil {
		return fmt.Errorf("failed to start service: %v", err)
	}

	// Limit the number of concurrent connections to the service.
	ln := netutil.LimitListener(l, serverMaxConnections)

	svc.logger.Info().Msgf("Vex service running on %s", l.Addr())

	err = svc.server.Serve(ln)
	if err == nil || err == http.ErrServerClosed {
		return nil
	}
	return fmt.Errorf("failed to serve service: %v", err)
}

// Stop gracefully shuts down the service. See [http.Server.Shutdown].
func (svc *Service) Stop(ctx context.Context) error {
	if err := svc.Validate(); err != nil {
		return err
	}

	err := svc.server.Shutdown(ctx)
	if err == nil || err == http.ErrServerClosed {
		return nil
	}
	return fmt.Errorf("failed to stop service: %v", err)
}
