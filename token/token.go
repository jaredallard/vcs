// Copyright (C) 2026 vcs contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

// Package token contains functions for getting an authenticated token
// from a user's machine for a given VCS provider.
package token

import (
	"context"
	"fmt"
	"time"

	"go.rgst.io/jaredallard/vcs/v3"
	"go.rgst.io/jaredallard/vcs/v3/token/internal/shared"
	"go.rgst.io/jaredallard/vcs/v3/token/providers/forgejo"
	"go.rgst.io/jaredallard/vcs/v3/token/providers/gitea"
	"go.rgst.io/jaredallard/vcs/v3/token/providers/github"
	"go.rgst.io/jaredallard/vcs/v3/token/providers/gitlab"
)

// defaultProviders contains all of the providers that are supported by
// this package by VCS provider.
var defaultProviders = map[vcs.Provider][]shared.Provider{
	vcs.ProviderGithub:  github.Providers,
	vcs.ProviderGitlab:  gitlab.Providers,
	vcs.ProviderGitea:   gitea.Providers,
	vcs.ProviderForgejo: forgejo.Providers,
}

// Token is a VCS token that can be used for API access. Defined here to
// allow for easy access to the type.
type Token = shared.Token

// Provider is an interface implemented by providers that return a
// [Token].
type Provider = shared.Provider

// Options contains options for the [Fetch] function.
type Options struct {
	// AllowUnauthenticated allows for an empty token to be returned if
	// no token is found. Defaults to false.
	AllowUnauthenticated bool

	// UseGlobalCache allows for the use of a global cache for tokens. If
	// set to true, the token will be cached globally (all instances of
	// this library). Otherwise, the token will always be fetched.
	//
	// Defaults to true.
	//
	// Note: When using [shared.Token], the value will never change.
	// Caching refers only to function calls provided by this package
	// (e.g., [Fetch]).
	UseGlobalCache *bool
}

// RegisterProvider sets the provided [Provider] implementations for the
// provided [vcs.Provider].
//
// Note: This REPLACES the existing providers, it is not additive. Also
// note that this mutates the global provider list and thus has
// side-effects. It should normally be used only when there is a custom
// provider required or potentially for test/mocking use-cases.
func RegisterProvider(vcsp vcs.Provider, p []Provider) error {
	defaultProviders[vcsp] = p
	return nil // for future usage
}

// applyDefaults applies defaults to the provided options.
func applyDefaults(opts *Options) {
	// If UseGlobalCache is not set, default to true.
	if opts.UseGlobalCache == nil {
		opts.UseGlobalCache = new(true)
	}
}

// Fetch returns a valid token from one of the configured credential
// providers. If no token is found, ErrNoToken is returned.
func Fetch(_ context.Context, vcsp vcs.Provider, opts *Options) (*shared.Token, error) {
	if _, ok := defaultProviders[vcsp]; !ok {
		return nil, fmt.Errorf("unknown VCS provider %q", vcsp)
	}

	if opts == nil {
		opts = &Options{}
	}

	applyDefaults(opts)

	if *opts.UseGlobalCache {
		t, ok := cache.Get(vcsp)
		if ok {
			return t.Clone(), nil
		}
	}

	var token *shared.Token
	var errs []error
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
		if !opts.AllowUnauthenticated {
			return nil, ErrNoToken(errs)
		}

		// Set an empty token since we're allowing unauthenticated access.
		token = &shared.Token{}
	}

	// Set when the token was fetched and store it in the cache for
	// possibly other calls to use.
	token.FetchedAt = time.Now()
	cache.Set(vcsp, token)

	return token, nil
}
