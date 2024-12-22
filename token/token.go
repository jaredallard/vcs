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

// Package token contains functions for getting an authenticated token
// from a user's machine for a given VCS provider.
package token

import (
	"context"
	"errors"
	"fmt"

	"github.com/jaredallard/vcs"
	"github.com/jaredallard/vcs/token/internal/github"
	"github.com/jaredallard/vcs/token/internal/gitlab"
	"github.com/jaredallard/vcs/token/internal/shared"
)

// defaultProviders contains all of the providers that are supported by
// this package by VCS provider.
var defaultProviders = map[vcs.Provider][]shared.Provider{
	vcs.ProviderGithub: github.Providers,
	vcs.ProviderGitlab: gitlab.Providers,
}

// Token is a VCS token that can be used for API access. Defined here to
// allow for easy access to the type.
type Token = shared.Token

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

// Fetch returns a valid token from one of the configured credential
// providers. If no token is found, ErrNoToken is returned.
//
// If allowUnauthenticated is true, then an empty token is returned if
// no token is found.
func Fetch(_ context.Context, vcsp vcs.Provider, allowUnauthenticated bool) (*shared.Token, error) {
	if _, ok := defaultProviders[vcsp]; !ok {
		return nil, fmt.Errorf("unknown VCS provider %q", vcsp)
	}

	var token *shared.Token
	errs := []error{}
	for _, p := range defaultProviders[vcsp] {
		var err error

		token, err = p.Token()
		if err != nil {
			errs = append(errs, err)
			continue
		}

		// Got a token, break out of the loop.
		if token != nil {
			break
		}
	}
	if token == nil {
		if allowUnauthenticated {
			return &shared.Token{}, nil
		}

		return nil, ErrNoToken(errs)
	}
	return token, nil
}
