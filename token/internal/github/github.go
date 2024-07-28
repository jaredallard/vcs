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

// Package github implements [shared.Provider] for Github.
package github

import (
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
	token, err := cmd.Output()
	if err != nil {
		return nil, execerr.From(cmd, err)
	}

	return &shared.Token{Value: strings.TrimSpace(string(token))}, nil
}
