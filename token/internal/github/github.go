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

// Package github implements [shared.Provider] for Github.
package github

import (
	"fmt"
	"strings"

	"github.com/jaredallard/cmdexec"
	"github.com/jaredallard/vcs/internal/execerr"
	"github.com/jaredallard/vcs/token/internal/shared"
)

// Providers is a list of providers that can be used to retrieve a
// token for Github.
var Providers = []shared.Provider{
	&shared.EnvProvider{EnvVars: []shared.EnvVar{{Name: "GITHUB_TOKEN"}, {Name: "GH_TOKEN"}}},
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
		return nil, execerr.From(cmd, err)
	}

	token := strings.TrimSpace(string(b))
	if token == "" {
		return nil, fmt.Errorf("no token returned from 'gh auth token'")
	}

	return &shared.Token{Value: token}, nil
}
