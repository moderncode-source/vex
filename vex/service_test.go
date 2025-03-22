// SPDX-FileCopyrightText: 2025 The Vex Authors.
//
// SPDX-License-Identifier: Apache-2.0 OR MIT
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. You may not use this file except in compliance with the
// terms of those licenses.

// Vex service tests.

package vex_test

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/moderncode-source/vex-svc/vex"
)

func TestNew(t *testing.T) {
	const addr = ":8080"
	const url = "http://localhost" + addr

	svc := vex.New(addr)
	var wg sync.WaitGroup

	// Give the entire operation with the service to complete within 1 second.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	defer func() {
		if err := svc.Stop(ctx); err != nil {
			t.Logf("Failed to stop service: %s", err)
		}
		wg.Wait()
	}()

	wg.Add(1)
	go func() {
		if err := svc.Start(); err != nil {
			t.Logf("Server exited with error: %s", err)
		}
		wg.Done()
	}()

	var err error
	var req *http.Request
	cli := &http.Client{}

	do := func(method, url string, wantStatus int) {
		t.Logf("Testing request to %s %s", method, url)

		req, err = http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %s", err)
		}

		resp, err := cli.Do(req)
		if err != nil {
			t.Fatalf("Request returned error: %s", err)
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Logf("Failed to close response body: %s", err)
			}
		}()

		if resp.StatusCode != wantStatus {
			t.Fatalf("Expected response code %d, got: %d", wantStatus, resp.StatusCode)
		}
	}

	// Default Vex service should at least respond to liveness and
	// readiness probe endpoints with 200 OK.
	do(http.MethodGet, url+vex.HealthEndpoint, http.StatusOK)
	do(http.MethodGet, url+vex.ReadyEndpoint, http.StatusOK)
}
