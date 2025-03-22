// SPDX-FileCopyrightText: 2025 The Vex Authors.
//
// SPDX-License-Identifier: Apache-2.0 OR MIT
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. You may not use this file except in compliance with the
// terms of those licenses.

// Vex HTTP request multiplexer. See [http.ServeMux].
// Request handlers are registered here.

package vex

import "net/http"

// ServiceMux is the default [http.ServeMux] used by [Service] with
// all handle functions automatically registered.
var ServiceMux = &serviceMux

var serviceMux http.ServeMux

const (
	// HealthEndpoint is an endpoint pattern that matches request of
	// any type to /healthz to [HealthHandler].
	HealthEndpoint = "/healthz"

	// ReadyEndpoint is an endpoint pattern that matches request of
	// any type to /v1/sys/ready to [ReadyHandler].
	ReadyEndpoint = "/v1/sys/ready"

	// PostQueueEndpoint is an endpoint pattern that matches POST requests to
	// /v1/queue to [PostQueueHandler] handler.
	PostQueueEndpoint = "POST /v1/queue/"

	// GetQueueEndpoint is an endpoint pattern that matches GET requests to
	// /v1/queue to [GetQueueHandler] handler.
	GetQueueEndpoint = "GET /v1/queue/"
)
