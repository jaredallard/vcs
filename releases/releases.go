// Copyright (C) 2026 vcs contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

// Package releases implements functions for interacting with 'Releases'
// provided by VCS providers. The Release terminology largely comes from
// Github and can be thought of as versioned artifacts that correspond
// to a Git tag.
package releases

import (
	"context"
	"fmt"
	"io"
	"io/fs"

	"go.rgst.io/jaredallard/vcs/v2"
	"go.rgst.io/jaredallard/vcs/v2/releases/gitea"
	"go.rgst.io/jaredallard/vcs/v2/releases/github"
	"go.rgst.io/jaredallard/vcs/v2/releases/gitlab"
	"go.rgst.io/jaredallard/vcs/v2/releases/internal/opts"
	"go.rgst.io/jaredallard/vcs/v2/token"
)

// fetchers is a map of VCS provider to their respective fetcher.
var fetchers = map[vcs.Provider]opts.Fetcher{
	vcs.ProviderGithub:  &github.Fetcher{},
	vcs.ProviderGitlab:  &gitlab.Fetcher{},
	vcs.ProviderGitea:   &gitea.Fetcher{},
	vcs.ProviderForgejo: &gitea.Fetcher{},
}

// GetReleaseNoteOptions is an alias for [opts.GetReleaseNoteOptions].
type GetReleaseNoteOptions = opts.GetReleaseNoteOptions

// FetchOptions is an alias for [opts.FetchOptions].
type FetchOptions = opts.FetchOptions

// Client contains configuration for fetching releases from various VCS
// providers.
type Client struct{}

// Fetch fetches a release from a VCS provider and returns an asset
// from it as an io.ReadCloser. This must be closed to close the
// underlying HTTP request.
//
//nolint:gocritic // Why: rc, name, size, error
func Fetch(ctx context.Context, opts *FetchOptions) (io.ReadCloser, fs.FileInfo, error) {
	if opts == nil {
		return nil, nil, fmt.Errorf("opts is nil")
	}

	if opts.RepoURL == "" {
		return nil, nil, fmt.Errorf("repo url is required")
	}

	if opts.Tag == "" {
		return nil, nil, fmt.Errorf("tag is required")
	}

	vcsp, err := vcs.ProviderFromURL(opts.RepoURL, opts.Overrides)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get VCS provider from URL: %w", err)
	}

	token, err := token.Fetch(ctx, vcsp, true)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch token: %w", err)
	}

	if fetcher, ok := fetchers[vcsp]; ok {
		return fetcher.Fetch(ctx, token, opts)
	}

	return nil, nil, fmt.Errorf("unknown VCS provider %s", vcsp)
}

// GetReleaseNotes fetches the release notes of a release from a VCS provider.
func GetReleaseNotes(ctx context.Context, opt *GetReleaseNoteOptions) (string, error) {
	if opt == nil {
		return "", fmt.Errorf("opts is nil")
	}

	if opt.RepoURL == "" {
		return "", fmt.Errorf("repo url is required")
	}

	if opt.Tag == "" {
		return "", fmt.Errorf("tag is required")
	}

	vcsp, err := vcs.ProviderFromURL(opt.RepoURL, opt.Overrides)
	if err != nil {
		return "", fmt.Errorf("failed to get VCS provider from URL: %w", err)
	}

	t, err := token.Fetch(ctx, vcsp, true)
	if err != nil {
		return "", fmt.Errorf("failed to fetch token: %w", err)
	}

	if fetcher, ok := fetchers[vcsp]; ok {
		return fetcher.GetReleaseNotes(ctx, t, opt)
	}

	return "", fmt.Errorf("unknown VCS provider %s", vcsp)
}
