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
				Value: token,
				Type:  env.Type,
			}, nil
		}
	}

	return nil, fmt.Errorf("no token found in environment variables: %v", p.EnvVars)
}
