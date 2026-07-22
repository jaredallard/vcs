// Copyright (C) 2026 vcs contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

// Package github implements [shared.Provider] for Github.
package github

import (
	"fmt"
	"strings"

	"go.rgst.io/jaredallard/cmdexec/v2"
	"go.rgst.io/jaredallard/vcs/v3/internal/execerr"
	"go.rgst.io/jaredallard/vcs/v3/token/internal/shared"
	"go.rgst.io/jaredallard/vcs/v3/token/providers/generic"
)

// Providers is a list of providers that can be used to retrieve a
// token for Github.
var Providers = []shared.Provider{
	&generic.EnvProvider{EnvVars: []generic.EnvVar{{Name: "GITHUB_TOKEN"}, {Name: "GH_TOKEN"}}},
	&GHProvider{},
}

// GHProvider implements the [token.Provider] interface using the Github
// CLI to retrieve a token.
type GHProvider struct{}

// Token returns a valid token or an error if no token is found.
func (p *GHProvider) Token() (*shared.Token, error) {
	cmd := cmdexec.Command("gh", "auth", "token")
	b, err := cmd.Output()
	if err != nil {
		return nil, execerr.From(err)
	}

	token := strings.TrimSpace(string(b))
	if token == "" {
		return nil, fmt.Errorf("no token returned from 'gh auth token'")
	}

	return &shared.Token{
		Source: "gh",
		Value:  token,
	}, nil
}
