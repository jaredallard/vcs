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

// Package gitlab contains Gitlab specific [token.Provider]s.
package gitlab

import (
	"strings"

	"github.com/jaredallard/cmdexec"
	"github.com/jaredallard/vcs/internal/execerr"
	"github.com/jaredallard/vcs/token/internal/shared"
)

// Contains the different types of tokens that can be retrieved.
const (
	TokenTypeJob = "job"
	TokenTypePAT = "pat"
)

// Providers is a list of providers that can be used to retrieve a
// token for Gitlab.
var Providers = []shared.Provider{
	envProvider(),
	&GlabProvider{},
}

// envProvider returns a [shared.EnvProvider] configured for Gitlab.
func envProvider() shared.Provider {
	return &shared.EnvProvider{EnvVars: []shared.EnvVar{
		{Name: "GITLAB_TOKEN"},
		{Name: "CI_JOB_TOKEN", Type: TokenTypeJob},
	}}
}

// GlabProvider implements the [token.Provider] interface using the
// Gitlab CLI (glab) to retrieve a token.
type GlabProvider struct{}

// Token returns a valid token or an error if no token is found.
func (p *GlabProvider) Token() (*shared.Token, error) {
	// determine the host from glab
	cmd := cmdexec.Command("glab", "config", "get", "-g", "host")
	out, err := cmd.Output()
	if err != nil {
		return nil, execerr.From(cmd, err)
	}
	host := strings.TrimSpace(string(out))

	cmd = cmdexec.Command("glab", "config", "get", "-g", "token", "-h", host)
	token, err := cmd.Output()
	if err != nil {
		return nil, execerr.From(cmd, err)
	}

	return &shared.Token{Value: strings.TrimSpace(string(token))}, nil
}
