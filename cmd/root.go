// SPDX-FileCopyrightText: 2025 The Vex Authors.
//
// SPDX-License-Identifier: Apache-2.0 OR MIT
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. You may not use this file except in compliance with the
// terms of those licenses.

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/moderncode-source/vex-svc/vex"
	"github.com/spf13/cobra"
)

const (
	verboseFlag = "verbose"
	addrFlag    = "addr"
)

func init() {
	rootCmd.Flags().BoolP(verboseFlag, "v", false, "run Vex with debug logs")
	rootCmd.Flags().String(addrFlag, ":8080", "the TCP network address for the Vex server to listen on")
}

var rootCmd = &cobra.Command{
	Use:   "vex",
	Short: "Run arbitrary code under isolated environments",
	Long: `Vex is a virtual execution micro-service that runs arbitrary code
in the cloud under controlled, isolated environments.
Documentation is available at https://github.com/moderncode-source/vex-svc`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		debug, err := cmd.Flags().GetBool(verboseFlag)
		if err != nil {
			return fmt.Errorf("could not retrieve %s flag value: %v", verboseFlag, err)
		}

		// TODO; replace fmt with a proper logger.
		if debug {
			fmt.Println("Welcome to Vex - a virtual execution micro-service")
			fmt.Println("Verbose output enabled")
		}

		addr, err := cmd.Flags().GetString(addrFlag)
		if err != nil {
			return fmt.Errorf("could not retrieve %s flag value: %v", addrFlag, err)
		}

		// TODO; replace fmt with a proper logger.
		// TODO: replace with svc.Start(), svc.Serve() for better validation.
		if debug {
			fmt.Printf("Will attempt to listen on: %s\n", addr)
			fmt.Println("Press CTRL+C to interrupt")
		}

		// Create new Vex service.
		svc := vex.New(addr)

		// TODO: might suit us better: [signal.NotifyContext].
		// If an interrupt signal is caught, gracefully shut down the
		// service and encapsulate any error it returns.
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)
		shutdownErr := make(chan error, 1)

		// TODO: add another chan to control this go-routine: wait for it to
		// exit before saying "goodbye" to the user.
		go func() {
			<-interrupt
			fmt.Println("CTRL+C detected, shutting down...")

			// TODO: pass a context with timeout.
			if err := svc.Stop(context.Background()); err != nil {
				shutdownErr <- err
			} else {
				close(shutdownErr)
			}
		}()

		if err := svc.Start(); err != nil {
			return fmt.Errorf("service error: %v", err)
		}

		if err := <-shutdownErr; err != nil {
			return fmt.Errorf("failed to shut down the service: %v", err)
		} else if debug {
			fmt.Println("Vex service was successfully shut down")
			fmt.Println("Goodbye!")
		}

		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}
