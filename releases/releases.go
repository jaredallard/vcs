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
	"strings"

	"github.com/jaredallard/vcs"
	"github.com/jaredallard/vcs/releases/github"
	"github.com/jaredallard/vcs/releases/internal/opts"
)

// fetchers is a map of VCS provider to their respective fetcher.
var fetchers = map[vcs.Provider]opts.Fetcher{
	vcs.ProviderGithub: &github.Fetcher{},
}

// FetchOptions is a set of options for Fetch
type FetchOptions struct {
	// RepoURL is the repository URL, it should be a valid
	// URL.
	RepoURL string

	// Tag is the tag of the release
	Tag string

	// AssetName is the name of the asset to fetch, globs are
	// supported.
	AssetName string

	// AssetNames is a list of asset names to fetch, the first
	// asset that matches will be returned. Globs are supported.
	AssetNames []string
}

// GetReleaseNoteOptions is a set of options for GetReleaseNotes
type GetReleaseNoteOptions struct {
	// RepoURL is the repository URL, it should be a valid
	// URL.
	RepoURL string

	// Tag is the tag of the release
	Tag string
}

// Client contains configuration for fetching releases from various VCS
// providers.
type Client struct{}

// Fetch fetches a release from a VCS provider and returns an asset
// from it as an io.ReadCloser. This must be closed to close the
// underlying HTTP request.
//
//nolint:gocritic // Why: rc, name, size, error
func Fetch(ctx context.Context, token string, opts *opts.FetchOptions) (io.ReadCloser, fs.FileInfo, error) {
	if opts == nil {
		return nil, nil, fmt.Errorf("opts is nil")
	}

	if opts.RepoURL == "" {
		return nil, nil, fmt.Errorf("repo url is required")
	}

	if opts.Tag == "" {
		return nil, nil, fmt.Errorf("tag is required")
	}

	if strings.Contains(opts.RepoURL, "github.com") {
		return fetchers[vcs.ProviderGithub].Fetch(ctx, token, opts)
	}

	return nil, nil, fmt.Errorf("unsupported fetch repo url: %s", opts.RepoURL)
}

// GetReleaseNotes fetches the release notes of a release from a VCS provider.
func GetReleaseNotes(ctx context.Context, token string, opts *opts.GetReleaseNoteOptions) (string, error) {
	if opts == nil {
		return "", fmt.Errorf("opts is nil")
	}

	if opts.RepoURL == "" {
		return "", fmt.Errorf("repo url is required")
	}

	if opts.Tag == "" {
		return "", fmt.Errorf("tag is required")
	}

	if strings.Contains(opts.RepoURL, "github.com") {
		return fetchers[vcs.ProviderGithub].GetReleaseNotes(ctx, token, opts)
	}

	return "", fmt.Errorf("unsupported get release notes repo url: %s", opts.RepoURL)
}
