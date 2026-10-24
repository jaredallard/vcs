// Copyright (C) 2026 vcs contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

// Description: Contains helpers for redacting credentials from URLs.

package git

import (
	"net/url"
	"strings"
)

// redactURL returns the provided URL with any credentials redacted to
// prevent them from being logged. For HTTP(S) URLs, all userinfo is
// redacted as tokens are commonly provided as just the username. For
// other schemes, only the password is redacted. If the URL can't be
// parsed or contains no userinfo, it is returned unchanged.
func redactURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil || u.User == nil {
		return rawURL
	}

	switch u.Scheme {
	case "http", "https":
		u.User = url.User("redacted")
	default:
		if _, ok := u.User.Password(); !ok {
			return rawURL
		}
		u.User = url.UserPassword(u.User.Username(), "redacted")
	}

	return u.String()
}

// redactedError wraps an error, replacing all occurrences of old with
// new in its message. The original error is still accessible through
// [errors.Unwrap].
type redactedError struct {
	err      error
	old, new string
}

// Error implements the error interface.
func (e *redactedError) Error() string {
	return strings.ReplaceAll(e.err.Error(), e.old, e.new)
}

// Unwrap returns the original error.
func (e *redactedError) Unwrap() error {
	return e.err
}

// redactErr returns err with any occurrences of rawURL in its message
// replaced by a redacted version of it, see [redactURL]. If err is nil
// or rawURL contains no credentials, err is returned unchanged.
func redactErr(err error, rawURL string) error {
	if err == nil {
		return nil
	}

	redacted := redactURL(rawURL)
	if redacted == rawURL {
		return err
	}

	return &redactedError{err: err, old: rawURL, new: redacted}
}
