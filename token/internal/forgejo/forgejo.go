// Copyright (C) 2026 vcs contributors
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

// Package forgejo contains Forgejo specific [token.Provider]s.
package forgejo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"go.rgst.io/jaredallard/cmdexec/v2"
	"go.rgst.io/jaredallard/vcs/v2/internal/execerr"
	"go.rgst.io/jaredallard/vcs/v2/token/internal/shared"
)

// signedInRegexp is used to parse the active user
var signedInRegexp = regexp.MustCompile(`currently signed in to (.*)@(.*)$`)

// Providers is a list of providers that can be used to retrieve a
// token for Forgejo.
var Providers = []shared.Provider{
	&shared.EnvProvider{EnvVars: []shared.EnvVar{{Name: "FORGEJO_TOKEM"}}},
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
