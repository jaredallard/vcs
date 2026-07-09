// Copyright (C) 2026 vcs contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package token

import (
	"sync"

	"go.rgst.io/jaredallard/vcs/v2"
	"go.rgst.io/jaredallard/vcs/v2/token/internal/shared"
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
