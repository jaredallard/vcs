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

// Package token contains functions for getting an authenticated token
// from a user's machine for a given VCS provider.
package token

import (
	"context"
	"errors"
	"fmt"

	"github.com/jaredallard/vcs"
	"github.com/jaredallard/vcs/token/github"
)

// defaultProviders contains all of the providers that are supported by
// this package by VCS provider.
var defaultProviders = map[vcs.Provider][]Provider{
	vcs.ProviderGithub: {
		&github.EnvProvider{},
		&github.GHProvider{},
	},
}

// ErrNoToken is returned when no token is found in the configured
// credential providers.
type ErrNoToken struct {
	errs []error
}

// Unwrap returns the errors that caused the ErrNoToken error.
func (e ErrNoToken) Unwrap() []error {
	return e.errs
}

// Error returns the error message for ErrNoToken.
func (e ErrNoToken) Error() string {
	return errors.Join(e.errs...).Error()
}

// Provider is an interface for VCS providers to implement to provide a
// token from a user's machine.
type Provider interface {
	// Token returns a valid token or an error if no token is found.
	Token() (string, error)
}

// Fetch returns a valid token from one of the configured credential
// providers. If no token is found, ErrNoToken is returned.
func Fetch(_ context.Context, vcsp vcs.Provider) (string, error) {
	if _, ok := defaultProviders[vcsp]; !ok {
		return "", fmt.Errorf("unknown VCS provider %q", vcsp)
	}

	token := ""
	errors := []error{}
	for _, p := range defaultProviders[vcsp] {
		var err error
		token, err = p.Token()
		if err != nil {
			errors = append(errors, err)
			continue
		}

		// Got a token, break out of the loop.
		if token != "" {
			break
		}
	}
	if token == "" {
		return "", ErrNoToken{errors}
	}
	return token, nil
}
