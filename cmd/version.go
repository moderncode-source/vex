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
	"fmt"

	"github.com/moderncode-source/vex-svc/vex"
	"github.com/spf13/cobra"
)

// Vex CLI major, minor, and patch version numbers.
const (
	VersionMajor = 0
	VersionMinor = 0
	VersionPatch = 1
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Vex",
	Args:  cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("Vex Core Version: v%d.%d.%d\n", vex.VersionMajor, vex.VersionMinor, vex.VersionPatch)
		fmt.Printf("Vex CLI Version: v%d.%d.%d\n", VersionMajor, VersionMinor, VersionPatch)
	},
}
