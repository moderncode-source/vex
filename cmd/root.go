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
	"sync"
	"syscall"

	"github.com/moderncode-source/vex-svc/vex"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
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
	SilenceUsage: true, // Do not print usage on error.
	RunE: func(cmd *cobra.Command, _ []string) error {
		// Get a logger for this command.
		var debug bool
		var err error

		cmdLogger := globalCmdLogger

		if debug, err = cmd.Flags().GetBool(verboseFlag); err != nil {
			return fmt.Errorf("could not retrieve %s flag value: %v", verboseFlag, err)
		} else if !debug {
			// https://github.com/rs/zerolog/tree/master#leveled-logging
			cmdLogger = cmdLogger.Level(zerolog.InfoLevel)
		}

		addr, err := cmd.Flags().GetString(addrFlag)
		if err != nil {
			return fmt.Errorf("could not retrieve %s flag value: %v", addrFlag, err)
		}

		cmdLogger.Info().Msg("Welcome to Vex - a virtual execution micro-service")

		// Create a thread-safe and fast logger for the service.
		w := diode.NewWriter(logOutput, diodeWriterSize, 0, func(missed int) {
			cmdLogger.Warn().Msgf("Service dropped %d logs", missed)
		})

		svcLogger := zerolog.New(w)
		if !debug {
			svcLogger = svcLogger.Level(zerolog.InfoLevel)
		}

		// Create a new Vex service.
		svc, err := vex.New(addr, &svcLogger)
		if err != nil {
			cmdLogger.Error().Err(err).Msg("Service creation error")
			return fmt.Errorf("service creation error: %v", err)
		}

		if err := svc.Validate(); err != nil {
			cmdLogger.Error().Err(err).Msg("Service error")
			return fmt.Errorf("service error: %v", err)
		}

		cmdLogger.Info().Msgf("Starting service process [%d] (Press CTRL+C to quit)", os.Getpid())

		// If an interrupt signal is caught, gracefully shut down the service
		// and encapsulate any error it returns.
		shutdownErr := make(chan error, 1)
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

		var wg sync.WaitGroup
		defer wg.Wait() // Wait for the go-routine below to exit.
		defer cancel()

		wg.Add(1)
		go func() {
			defer wg.Done()

			<-ctx.Done()
			cmdLogger.Info().Msg("Shutting down...")

			// TODO: pass a context with timeout.
			if err := svc.Stop(context.Background()); err != nil {
				shutdownErr <- err
			} else {
				close(shutdownErr)
			}
		}()

		if err := svc.Start(); err != nil {
			cmdLogger.Error().Err(err).Msg("Service start error")
			return fmt.Errorf("service start error: %v", err)
		}

		if err := <-shutdownErr; err != nil {
			cmdLogger.Error().Err(err).Msg("Service shutdown error")
			return fmt.Errorf("service shutdown error: %v", err)
		}

		cmdLogger.Info().Msg("Service shutdown complete")
		cmdLogger.Info().Msgf("Finished service process [%d]", os.Getpid())
		return nil
	},
}

// Execute matches ands runs the appropriate CLI command.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("vex CLI exited with error: %s", err)
	}
	return nil
}
