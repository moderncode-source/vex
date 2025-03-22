// SPDX-FileCopyrightText: 2025 The Vex Authors.
//
// SPDX-License-Identifier: Apache-2.0 OR MIT
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. You may not use this file except in compliance with the
// terms of those licenses.

// Vex HTTP request handlers tests.

package vex_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/moderncode-source/vex-svc/vex"
)

const mockURL = "http://localhost:8080"

func checkHandlerResponseCode(
	ctx context.Context, handler http.Handler, want int, method, url string,
) error {
	res := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(ctx, method, url, nil)

	handler.ServeHTTP(res, req)
	if res.Code != want {
		return fmt.Errorf("Expected response code %d, got %d", want, res.Code)
	}

	return nil
}

func TestHealthHandler(t *testing.T) {
	const url = mockURL + vex.HealthEndpoint

	ctx := context.Background()

	if err := checkHandlerResponseCode(
		ctx, http.HandlerFunc(vex.HealthHandler),
		http.StatusOK, http.MethodGet, url,
	); err != nil {
		t.Fatal(err)
	}

	if err := checkHandlerResponseCode(
		ctx, http.HandlerFunc(vex.HealthHandler),
		http.StatusMethodNotAllowed, http.MethodPost, url,
	); err != nil {
		t.Fatal(err)
	}
}

func TestReadyHandler(t *testing.T) {
	const url = mockURL + vex.ReadyEndpoint

	ctx := context.Background()

	if err := checkHandlerResponseCode(
		ctx, http.HandlerFunc(vex.ReadyHandler),
		http.StatusOK, http.MethodGet, url,
	); err != nil {
		t.Fatal(err)
	}

	if err := checkHandlerResponseCode(
		ctx, http.HandlerFunc(vex.ReadyHandler),
		http.StatusMethodNotAllowed, http.MethodPost, url,
	); err != nil {
		t.Fatal(err)
	}
}

// TODO: test [vex.PostQueueHandler], [vex.GetQueueHandler].
