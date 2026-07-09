// Copyright (C) 2026 vcs contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

// Package gitlab contains Gitlab specific [token.Provider]s.
package gitlab

import (
	"fmt"
	"strings"

	"go.rgst.io/jaredallard/cmdexec/v2"
	"go.rgst.io/jaredallard/vcs/v2/internal/execerr"
	"go.rgst.io/jaredallard/vcs/v2/token/internal/shared"
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
