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

// mockService is used to test request handler responses.
var mockService = vex.NewWithHandler("", nil, nil)

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

func TestHealthAndReadyHandlers(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		ctx         context.Context
		handler     http.Handler
		want        int
		method, url string
	}{
		{ctx, http.HandlerFunc(mockService.HealthHandler), http.StatusOK, http.MethodGet, mockURL + vex.HealthEndpoint},
		{ctx, http.HandlerFunc(mockService.HealthHandler), http.StatusMethodNotAllowed, http.MethodPost, mockURL + vex.HealthEndpoint},

		{ctx, http.HandlerFunc(mockService.ReadyHandler), http.StatusOK, http.MethodGet, mockURL + vex.ReadyEndpoint},
		{ctx, http.HandlerFunc(mockService.ReadyHandler), http.StatusMethodNotAllowed, http.MethodPost, mockURL + vex.ReadyEndpoint},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s %s", tt.method, tt.url)

		t.Run(testname, func(t *testing.T) {
			if err := checkHandlerResponseCode(
				tt.ctx, tt.handler, tt.want, tt.method, tt.url,
			); err != nil {
				t.Fatal(err)
			}
		})
	}
}

// TODO: test [vex.PostQueueHandler], [vex.GetQueueHandler].
