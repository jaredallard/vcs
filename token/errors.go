// Copyright (C) 2026 vcs contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package token

import "errors"

// ErrNoToken is returned when no token is found in the configured
// credential providers.
type ErrNoToken []error

// Unwrap returns the errors that caused the ErrNoToken error.
func (errs ErrNoToken) Unwrap() []error {
	return errs
}

// Error returns the error message for ErrNoToken.
func (errs ErrNoToken) Error() string {
	return errors.Join(errs...).Error()
}

// NewErrNoToken returns a new [ErrNoToken] from the provided errors.
// Equivalent to ErrNoToken(errs).
func NewErrNoToken(errs []error) ErrNoToken {
	return ErrNoToken(errs)
}
