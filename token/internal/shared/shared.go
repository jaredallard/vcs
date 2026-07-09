// Copyright (C) 2026 vcs contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

// Package shared contains shared provider implementations. Largely
// stored in this package to avoid circular dependencies.
package shared

import (
	"strings"
	"time"
)

// Token is a VCS token that can be used for API access.
//
// Do not use the 'shared.Token' type, instead use [token.Token] which
// is an alias to this type.
type Token struct {
	// FetchedAt is the time that the token was fetched at. This does not
	// need to be set by providers as it is set by the [token.Fetch]
	// function.
	FetchedAt time.Time

	// Value is the token value.
	Value string

	// Source is the source of the token, this is set depending on the
	// provider that provided the token (e.g., `gh` for the Github CLI).
	Source string

	// Type is the type of the token, this is set depending on the
	// provider that provided the token.
	Type string
}

// IsUnauthenticated returns true if the token is empty.
func (t *Token) IsUnauthenticated() bool {
	return t.Value == ""
}

// String returns a redacted version of the token to prevent accidental
// logging.
func (t *Token) String() string {
	// keep the first 4 characters of the token, redact the rest.
	if len(t.Value) > 4 {
		prefix := t.Value[:4]
		return prefix + strings.Repeat("*", len(t.Value)-4)
	}

	// otherwise return the full token, but this is probably an invalid
	// token.
	return t.Value
}

// Clone returns a deep clone of the token.
func (t *Token) Clone() *Token {
	return &Token{
		FetchedAt: t.FetchedAt,
		Source:    t.Source,
		Value:     t.Value,
		Type:      t.Type,
	}
}

// Provider is an interface for VCS providers to implement to provide a
// token from a user's machine.
type Provider interface {
	// Token returns a valid token or an error if no token is found.
	Token() (*Token, error)
}
