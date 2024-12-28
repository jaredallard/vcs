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
