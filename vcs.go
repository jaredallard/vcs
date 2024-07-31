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

// Package vcs contains constants for the VCS providers supported by
// the libraries provided.
package vcs

import (
	"fmt"
	"strings"
)

// Provider represents a VCS provider.
type Provider string

// Contains constants for provider enum values.
const (
	// ProviderGithub represents Github.
	ProviderGithub Provider = "github"

	// ProviderGitlab represents Gitlab.
	ProviderGitlab Provider = "gitlab"
)

// Override represents an override for a given URL passed to
// ProviderFromURL.
type Override struct {
	// URLBase is the base URL that this override should apply to.
	URLBase string

	// Provider is the provider to override to.
	Provider Provider
}

// ProviderFromURL returns the VCS provider from a URL.
func ProviderFromURL(url string, overrides []Override) (Provider, error) {
	// Check for overrides.
	for _, override := range overrides {
		if strings.HasPrefix(url, override.URLBase) {
			return override.Provider, nil
		}
	}

	// Otherwise, fallback to heuristics.
	switch {
	case strings.Contains(url, "github.com"):
		return ProviderGithub, nil
	case strings.Contains(url, "gitlab.com"):
		return ProviderGitlab, nil
	case strings.Contains(url, "gitlab."):
		// Support gitlab.xyz addresses.
		return ProviderGitlab, nil
	default:
		return "", fmt.Errorf("unknown VCS provider for URL: %s", url)
	}
}
