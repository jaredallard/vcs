// Copyright (C) 2024 Jared Allard <jaredallard@users.noreply.github.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

// Package shared contains shared provider implementations. Largely
// stored in this package to avoid circular dependencies.
package shared

import (
	"strings"
)

// Token is a VCS token that can be used for API access.
//
// Do not use the 'shared.Token' type, instead use [token.Token] which
// is an alias to this type.
type Token struct {
	// Value is the token value.
	Value string

	// Type is the type of the token, this is set depending on the
	// provider that provided the token.
	Type string
}

// String returns a redacted version of the token to prevent accidental
// logging.
func (t *Token) String() string {
	// keep the first 4 characters of the token, redact the rest.
	if len(t.Value) > 4 {
		prefix := t.Value[:4]
		return prefix + strings.Repeat("*", len(t.Value)-4)
	}

	// other wise return the full token, but this is probably an invalid
	// token.
	return t.Value
}

// Provider is an interface for VCS providers to implement to provide a
// token from a user's machine.
type Provider interface {
	// Token returns a valid token or an error if no token is found.
	Token() (*Token, error)
}
