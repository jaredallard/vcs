// Copyright (C) 2024 vcs contributors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this program. If not, see
// <https://www.gnu.org/licenses/>.
//
// SPDX-License-Identifier: LGPL-3.0

// Package execerr contains an error handler for errors returned by the
// exec library. When an exec.ExitError is returned, it's formatted into
// an error with details about the command that was executed and it's
// output as present on the exec.ExitError.
package execerr

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/jaredallard/cmdexec"
)

// From attempts to parse the provided error as an exec.ExitError. If
// it's not an exec.ExitError, the original error is returned unchanged.
//
// Otherwise, a new error is returned with details about the command
// that was executed and it's output as present on the exec.ExitError.
func From(_ cmdexec.Cmd, err error) error {
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
