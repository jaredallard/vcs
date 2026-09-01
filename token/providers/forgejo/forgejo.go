// Copyright (C) 2026 vcs contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

// Package forgejo contains Forgejo specific [token.Provider]s.
package forgejo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"go.rgst.io/jaredallard/cmdexec/v2"
	"go.rgst.io/jaredallard/vcs/v3/internal/execerr"
	"go.rgst.io/jaredallard/vcs/v3/token/internal/shared"
	"go.rgst.io/jaredallard/vcs/v3/token/providers/generic"
)

// signedInRegexp is used to parse the active user
var signedInRegexp = regexp.MustCompile(`currently signed in to (.*)@(.*)$`)

// Providers is a list of providers that can be used to retrieve a
// token for Forgejo.
var Providers = []shared.Provider{
	&generic.EnvProvider{EnvVars: []generic.EnvVar{{Name: "FORGEJO_TOKEN"}}},
	&FJProvider{},
}

// FJProvider implements the [token.Provider] interface using the
// Forgejo CLI to retrieve a token.
type FJProvider struct{}

// FJKeys represents ~/.local/share/forgejo-cli/keys.json
type FJKeys struct {
	Hosts      map[string]FJKeysHost `json:"hosts"`
	Aliases    any                   `json:"aliases"`
	DefaultSSH []any                 `json:"default_ssh"`
}

// FJKeysHost is a host within [FJKeys].
type FJKeysHost struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    []int  `json:"expires_at"`
}

// Token returns a valid token or an error if no token is found.
func (p *FJProvider) Token() (*shared.Token, error) {
	cmd := cmdexec.Command("fj", "whoami")
	b, err := cmd.Output()
	if err != nil {
		return nil, execerr.From(err)
	}

	matches := signedInRegexp.FindStringSubmatch(string(b))
	if len(matches) != 3 {
		return nil, fmt.Errorf("failed to get current user from %q (got: %s)", cmd.String(), string(b))
	}

	host := matches[2]

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user's home directory: %w", err)
	}

	keysPath := filepath.Join(homeDir, ".local", "share", "forgejo-cli", "keys.json")
	f, err := os.Open(keysPath) //nolint:gosec // Why: intentional
	if err != nil {
		return nil, fmt.Errorf("failed to read forgejo keys.json at %q: %w", keysPath, err)
	}
	defer f.Close() //nolint:errcheck // Why: acceptable

	var keys FJKeys
	if err := json.NewDecoder(f).Decode(&keys); err != nil {
		return nil, fmt.Errorf("failed to decodes forgejo keys.json at %q: %w", keysPath, err)
	}

	if _, ok := keys.Hosts[host]; !ok {
		return nil, fmt.Errorf("no keys for host %q", keysPath)
	}

	return &shared.Token{
		Source: "fj",
		Value:  keys.Hosts[host].Token,
	}, nil
}
