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

package token

import (
	"sync"

	"github.com/jaredallard/vcs"
	"github.com/jaredallard/vcs/token/internal/shared"
)

// tokenCache is a cache of tokens that have been fetched from the
// user's machine.
type tokenCache struct {
	// tokensMu is a mutex to protect the tokens map.
	tokensMu sync.RWMutex

	// tokens is a map of VCS provider to their respective token.
	tokens map[vcs.Provider]*shared.Token
}

// Get returns a token from the cache if it exists.
func (c *tokenCache) Get(provider vcs.Provider) (*shared.Token, bool) {
	c.tokensMu.RLock()
	defer c.tokensMu.RUnlock()

	t, ok := c.tokens[provider]
	return t, ok
}

// Set sets a token in the cache.
func (c *tokenCache) Set(provider vcs.Provider, token *shared.Token) {
	c.tokensMu.Lock()
	defer c.tokensMu.Unlock()

	c.tokens[provider] = token
}

// cache is the global token cache.
var cache = &tokenCache{tokens: make(map[vcs.Provider]*shared.Token)}
