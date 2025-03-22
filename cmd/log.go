// SPDX-FileCopyrightText: 2025 The Vex Authors.
//
// SPDX-License-Identifier: Apache-2.0 OR MIT
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. You may not use this file except in compliance with the
// terms of those licenses.

// Vex CLI logging configuration.

package cmd

import (
	"os"

	"github.com/rs/zerolog"
)

// diodeWriterSize determines the buffer size of a diode ring buffer used for
// a thread-safe logging writer. See [diode.NewWriter] for details. If the
// value is too small, stderr will contain logs mentioning this.
const diodeWriterSize = 100

var logOutput = os.Stderr

// Human-friendly loggers used by Vex CLI commands to log output in a pretty
// format. The Vex service itself uses a thread-safe and fast logger instead.
var (
	globalCmdLogger = zerolog.New(zerolog.ConsoleWriter{Out: logOutput}).With().Timestamp().Logger()
)

func init() {
	// Set logging time-stamp format globally. This affects Vex service
	// logs if CLI passes a zerolog logger to it.
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}
