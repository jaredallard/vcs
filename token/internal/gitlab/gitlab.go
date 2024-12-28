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

// Package gitlab contains Gitlab specific [token.Provider]s.
package gitlab

import (
	"fmt"
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
	b, err := cmd.Output()
	if err != nil {
		return nil, execerr.From(err)
	}
	host := strings.TrimSpace(string(b))

	cmd = cmdexec.Command("glab", "config", "get", "-g", "token", "-h", host)
	b, err = cmd.Output()
	if err != nil {
		return nil, execerr.From(err)
	}

	token := strings.TrimSpace(string(b))
	if token == "" {
		return nil, fmt.Errorf("no token returned")
	}

	return &shared.Token{
		Source: "glab",
		Value:  token,
	}, nil
}
