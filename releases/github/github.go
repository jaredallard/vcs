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

// Package github implements [opts.Fetcher] for Github releases.
package github

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	gogithub "github.com/google/go-github/v78/github"
	"github.com/jaredallard/vcs/internal/fileinfo"
	"github.com/jaredallard/vcs/releases/internal/opts"
	"github.com/jaredallard/vcs/token"
	"golang.org/x/oauth2"
)

// _ is a compile-time assertion that Fetcher implements the
// [opts.Fetcher] interface.
var _ opts.Fetcher = &Fetcher{}

// Fetcher implements the [releases.Fetcher] interface for Github releases.
type Fetcher struct{}

// assetToFileInfo creates a type that satisfies [os.FileInfo] from the
// given [gogithub.ReleaseAsset].
func assetToFileInfo(a *gogithub.ReleaseAsset) os.FileInfo {
	modTime := a.UpdatedAt.Time
	if modTime.IsZero() {
		modTime = a.CreatedAt.Time
	}

	return fileinfo.New(a.GetName(), int64(a.GetSize()), modTime, a)
}

// getOrgRepoFromURL returns the org and repo from a URL:
//
// Example: https://github.com/rgst-io/stencil
func getOrgRepoFromURL(urlStr string) (owner, repo string, err error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", "", err
	}

	// /rgst-io/stencil -> ["", "rgst-io", "stencil"]
	spl := strings.Split(u.Path, "/")
	if len(spl) != 3 {
		return "", "", fmt.Errorf("invalid Github URL: %s", urlStr)
	}
	return spl[1], spl[2], nil
}

// createClient creates a Github client
func (f *Fetcher) createClient(ctx context.Context, t *token.Token) *gogithub.Client {
	httpClient := http.DefaultClient
	if !t.IsUnauthenticated() {
		httpClient = oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: t.Value}))
	}
	return gogithub.NewClient(httpClient)
}

// GetReleaseNotes returns the release notes for a given tag
func (f *Fetcher) GetReleaseNotes(ctx context.Context, t *token.Token, opt *opts.GetReleaseNoteOptions) (string, error) {
	gh := f.createClient(ctx, t)
	friendlyRepo := strings.TrimPrefix(opt.RepoURL, "https://")

	org, repo, err := getOrgRepoFromURL(opt.RepoURL)
	if err != nil {
		return "", err
	}

	rel, _, err := gh.Repositories.GetReleaseByTag(ctx, org, repo, opt.Tag)
	if err != nil {
		return "", fmt.Errorf("failed to get release for %s@%s: %w", friendlyRepo, opt.Tag, err)
	}

	return rel.GetBody(), nil
}

// Fetch fetches a release from a github repository and the underlying
// release asset.
func (f *Fetcher) Fetch(ctx context.Context, t *token.Token, opt *opts.FetchOptions) (io.ReadCloser, os.FileInfo, error) {
	gh := f.createClient(ctx, t)

	friendlyRepo := strings.TrimPrefix(opt.RepoURL, "https://")

	org, repo, err := getOrgRepoFromURL(opt.RepoURL)
	if err != nil {
		return nil, nil, err
	}

	rel, _, err := gh.Repositories.GetReleaseByTag(ctx, org, repo, opt.Tag)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get release for %s@%s: %w", friendlyRepo, opt.Tag, err)
	}

	// copy the assetNames slice, and append the assetName if it is not
	// empty
	validAssets := append([]string{}, opt.AssetNames...)
	if opt.AssetName != "" {
		validAssets = append(validAssets, opt.AssetName)
	}

	// Find an asset that matches the provided asset names
	var a *gogithub.ReleaseAsset
	for _, asset := range rel.Assets {
		for _, assetName := range validAssets {
			matched := false

			// attempt to use glob first, if that errors then fall back to
			// straight strings comparison
			if match, err := filepath.Match(assetName, asset.GetName()); err == nil {
				matched = match
			} else if assetName == asset.GetName() {
				matched = true
			}

			if matched {
				a = asset
				break
			}
		}
	}
	if a == nil {
		return nil, nil,
			fmt.Errorf("failed to find asset %v in release %s@%s", validAssets, friendlyRepo, opt.Tag)
	}

	// The second return value is a redirectURL, but by passing
	// http.DefaultClient we shouldn't have to handle it.
	rc, _, err := gh.Repositories.DownloadReleaseAsset(ctx, org, repo, a.GetID(), http.DefaultClient)
	if err != nil {
		return nil, nil,
			fmt.Errorf("failed to download asset %s from release %s@%s: %w", a.GetName(), friendlyRepo, opt.Tag, err)
	}

	return rc, assetToFileInfo(a), nil
}
