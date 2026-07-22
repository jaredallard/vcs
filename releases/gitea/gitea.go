// Copyright (C) 2026 vcs contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

// Package gitea implements [opts.Fetcher] for Gitea/Forgjo releases.
package gitea

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"code.gitea.io/sdk/gitea"
	"go.rgst.io/jaredallard/vcs/v3/internal/fileinfo"
	"go.rgst.io/jaredallard/vcs/v3/releases/internal/opts"
	"go.rgst.io/jaredallard/vcs/v3/token"
)

// _ is a compile-time assertion that Fetcher implements the
// [opts.Fetcher] interface.
var _ opts.Fetcher = &Fetcher{}

// Fetcher implements the [releases.Fetcher] interface for Gitea releases.
type Fetcher struct{}

// assetToFileInfo creates a type that satisfies [os.FileInfo] from the
// given [gitea.Attachment].
func assetToFileInfo(a *gitea.Attachment) os.FileInfo {
	return fileinfo.New(a.Name, a.Size, a.Created, a)
}

// createClient creates a Gitea client
func (f *Fetcher) createClient(repoURL string, t *token.Token) (*gitea.Client, error) {
	u, err := url.Parse(repoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse repo url %q as a URL: %w", repoURL, err)
	}

	baseURL := u.Scheme + "://" + u.Host
	if t.IsUnauthenticated() {
		return gitea.NewClient(baseURL)
	}

	client, err := gitea.NewClient(u.Scheme+"://"+u.Host, gitea.SetBasicAuth("token", t.Value))
	return client, err
}

// getOrgRepoFromURL returns the org and repo from a URL:
//
// Example: https://git.rgst.io/rgst-io/stencil
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

// GetReleaseNotes returns the release notes for a given tag
func (f *Fetcher) GetReleaseNotes(_ context.Context, t *token.Token, opt *opts.GetReleaseNoteOptions) (string, error) {
	gsdk, err := f.createClient(opt.RepoURL, t)
	if err != nil {
		return "", err
	}
	friendlyRepo := strings.TrimPrefix(opt.RepoURL, "https://")

	owner, repo, err := getOrgRepoFromURL(opt.RepoURL)
	if err != nil {
		return "", err
	}

	rel, _, err := gsdk.GetReleaseByTag(owner, repo, opt.Tag)
	if err != nil {
		return "", fmt.Errorf("failed to get release for %s@%s: %w", friendlyRepo, opt.Tag, err)
	}
	return rel.Note, nil
}

// Fetch fetches a release from a github repository and the underlying
// release asset.
func (f *Fetcher) Fetch(_ context.Context, t *token.Token, opt *opts.FetchOptions) (io.ReadCloser, os.FileInfo, error) {
	gsdk, err := f.createClient(opt.RepoURL, t)
	if err != nil {
		return nil, nil, err
	}
	friendlyRepo := strings.TrimPrefix(opt.RepoURL, "https://")

	owner, repo, err := getOrgRepoFromURL(opt.RepoURL)
	if err != nil {
		return nil, nil, err
	}

	rel, _, err := gsdk.GetReleaseByTag(owner, repo, opt.Tag)
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
	var a *gitea.Attachment
	for _, asset := range rel.Attachments {
		for _, assetName := range validAssets {
			matched := false

			// attempt to use glob first, if that errors then fall back to
			// straight strings comparison
			if match, err := filepath.Match(assetName, asset.Name); err == nil {
				matched = match
			} else if assetName == asset.Name {
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

	// Download the asset
	req, err := http.NewRequest(http.MethodGet, a.DownloadURL, http.NoBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request to download asset: %w", err)
	}

	if !t.IsUnauthenticated() {
		req.Header.Set("Authorization", "token "+t.Value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil,
			fmt.Errorf("failed to download asset %s from release %s@%s: %w", a.Name, friendlyRepo, opt.Tag, err)
	}
	return resp.Body, assetToFileInfo(a), nil
}
