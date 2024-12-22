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

	"github.com/jaredallard/vcs"
	"github.com/jaredallard/vcs/releases/github"
	"github.com/jaredallard/vcs/releases/gitlab"
	"github.com/jaredallard/vcs/releases/internal/opts"
	"github.com/jaredallard/vcs/token"
)

// fetchers is a map of VCS provider to their respective fetcher.
var fetchers = map[vcs.Provider]opts.Fetcher{
	vcs.ProviderGithub: &github.Fetcher{},
	vcs.ProviderGitlab: &gitlab.Fetcher{},
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
