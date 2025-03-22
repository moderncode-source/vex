// SPDX-FileCopyrightText: 2025 The Vex Authors.
//
// SPDX-License-Identifier: Apache-2.0 OR MIT
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. You may not use this file except in compliance with the
// terms of those licenses.

// Vex submission primitives.

package vex

// TODO: consider renaming this file.

// Submission type represents a code execution request data sent by the client.
// Submissions are queued and later scheduled to execute.
// TODO: replace "queue" in the description here with a type reference.
type Submission struct {
	// ID is a unique identifier for this submission.
	ID int64 `json:"id"`

	// Timestamp stores the time in RFC3339 format at which this submission was
	// created by the client. Used to monitor and estimate an average total
	// time it takes to execute a submission with respect to current traffic.
	// TODO: is this description accurate?
	// TODO: consider UNIX timestamp: smaller and will be processed faster.
	//       **OR** do two fields: one for RFC3339 and another for UNIX.
	//       Then let your service config set what is the preferred one.
	Timestamp string `json:"timestamp"`

	// TODO: add more fields.
}
