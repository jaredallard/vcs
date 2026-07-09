// Copyright (C) 2026 vcs contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package shared

import (
	"fmt"
	"os"
)

// EnvVar is a struct that represents an environment variable that can
// contain a VCS token.
type EnvVar struct {
	// Name is the name of the environment variable.
	Name string

	// Type is an optional field that denotes what type of token this.
	Type string
}

// EnvProvider implements the [token.Provider] interface using the
// environment variables to retrieve a token.
type EnvProvider struct {
	// EnvVars is a list of environment variables to check for a token.
	EnvVars []EnvVar
}

// Token returns a valid token or an error if no token is found.
func (p *EnvProvider) Token() (*Token, error) {
	for _, env := range p.EnvVars {
		if token := os.Getenv(env.Name); token != "" {
			return &Token{
				Value:  token,
				Source: fmt.Sprintf("environment variable (%s)", env.Name),
				Type:   env.Type,
			}, nil
		}
	}

	return nil, fmt.Errorf("no token found in environment variables: %v", p.EnvVars)
}
