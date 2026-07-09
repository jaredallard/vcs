// Copyright (C) 2026 vcs contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

// Package execerr contains an error handler for errors returned by the
// exec library. When an exec.ExitError is returned, it's formatted into
// an error with details about the command that was executed and it's
// output as present on the exec.ExitError.
package execerr

import (
	"errors"
	"fmt"
	"os/exec"
)

// From attempts to parse the provided error as an exec.ExitError. If
// it's not an exec.ExitError, the original error is returned unchanged.
//
// Otherwise, a new error is returned with details about the command
// that was executed and it's output as present on the exec.ExitError.
func From(err error) error {
	if err == nil {
		return nil
	}

	var execErr *exec.ExitError
	if !errors.As(err, &execErr) {
		return fmt.Errorf("exec failed (not *exec.ExitError): %w", err)
	}

	stderr := string(execErr.Stderr)
	if stderr == "" {
		stderr = "[no stderr]"
	}
	return fmt.Errorf("exec failed (%w): %s", execErr, stderr)
}
