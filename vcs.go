// Copyright (C) 2026 vcs contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

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

	// ProviderGitea represents Gitea.
	ProviderGitea Provider = "gitea"

	// ProviderForgejo represents Forgejo.
	ProviderForgejo Provider = "forgejo"
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
	case strings.Contains(url, "codeberg.org"), strings.Contains(url, "git.rgst.io"):
		return ProviderForgejo, nil
	default:
		return "", fmt.Errorf("unknown VCS provider for URL: %s", url)
	}
}
