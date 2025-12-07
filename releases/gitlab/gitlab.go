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

// Package gitlab implements [opts.Fetcher] for Gitlab releases.
package gitlab

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jaredallard/vcs/internal/fileinfo"
	"github.com/jaredallard/vcs/releases/internal/opts"
	"github.com/jaredallard/vcs/token"
	gogitlab "gitlab.com/gitlab-org/api/client-go"
)

// _ is a compile-time assertion that Fetcher implements the
// [opts.Fetcher] interface.
var _ opts.Fetcher = &Fetcher{}

// Fetcher implements the [releases.Fetcher] interface for Gitlab releases.
type Fetcher struct{}

// assetToFileInfo creates a type that satisfies [os.FileInfo] from the
// given [gogitlab.ReleaseLink].
func assetToFileInfo(rl *gogitlab.ReleaseLink) os.FileInfo {
	return fileinfo.New(rl.Name, 0, time.Time{}, rl)
}

// createClient creates a Gitlab client
func (f *Fetcher) createClient(t *token.Token) (*gogitlab.Client, error) {
	if t.IsUnauthenticated() {
		return gogitlab.NewClient("")
	}

	var client *gogitlab.Client
	var err error
	switch t.Type {
	case "pat", "": // Default is PAT.
		client, err = gogitlab.NewClient(t.Value)
	case "job":
		client, err = gogitlab.NewJobClient(t.Value)
	default:
		return nil, fmt.Errorf("unknown token type %s", t.Type)
	}
	return client, err
}

// getPIDFromRepoURL returns the project ID from a given repository URL.
func (f *Fetcher) getPIDFromRepoURL(repoURL string, glab *gogitlab.Client) (int64, error) {
	u, err := url.Parse(repoURL)
	if err != nil {
		return 0, err
	}

	proj, _, err := glab.Projects.GetProject(strings.TrimPrefix(u.Path, "/"), nil)
	if err != nil {
		return 0, err
	}

	return proj.ID, nil
}

// GetReleaseNotes returns the release notes for a given tag
func (f *Fetcher) GetReleaseNotes(_ context.Context, t *token.Token, opt *opts.GetReleaseNoteOptions) (string, error) {
	glab, err := f.createClient(t)
	if err != nil {
		return "", err
	}

	friendlyRepo := strings.TrimPrefix(opt.RepoURL, "https://")
	pid, err := f.getPIDFromRepoURL(opt.RepoURL, glab)
	if err != nil {
		return "", err
	}

	rel, _, err := glab.Releases.GetRelease(pid, opt.Tag)
	if err != nil {
		return "", fmt.Errorf("failed to get release for %s@%s: %w", friendlyRepo, opt.Tag, err)
	}
	return rel.Description, nil
}

// Fetch fetches a release from a github repository and the underlying
// release asset.
func (f *Fetcher) Fetch(_ context.Context, t *token.Token, opt *opts.FetchOptions) (io.ReadCloser, os.FileInfo, error) {
	glab, err := f.createClient(t)
	if err != nil {
		return nil, nil, err
	}

	friendlyRepo := strings.TrimPrefix(opt.RepoURL, "https://")
	pid, err := f.getPIDFromRepoURL(opt.RepoURL, glab)
	if err != nil {
		return nil, nil, err
	}

	rel, _, err := glab.Releases.GetRelease(pid, opt.Tag)
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
	var rl *gogitlab.ReleaseLink
	for _, relLink := range rel.Assets.Links {
		for _, assetName := range validAssets {
			matched := false

			// attempt to use glob first, if that errors then fall back to
			// straight strings comparison
			if match, err := filepath.Match(assetName, relLink.Name); err == nil {
				matched = match
			} else if assetName == relLink.Name {
				matched = true
			}

			if matched {
				rl = relLink
				break
			}
		}
	}
	if rl == nil {
		return nil, nil,
			fmt.Errorf("failed to find asset %v in release %s@%s", validAssets, friendlyRepo, opt.Tag)
	}

	// Download the asset
	req, err := http.NewRequest(http.MethodGet, rl.DirectAssetURL, http.NoBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request to download asset: %w", err)
	}
	// TODO(jaredallard): Gitlab's auth system is awful, so job token
	// won't _just work_. We'll eventually need to support it.
	req.Header.Set("PRIVATE-TOKEN", t.Value)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil,
			fmt.Errorf("failed to download asset %s from release %s@%s: %w", rl.Name, friendlyRepo, opt.Tag, err)
	}
	return resp.Body, assetToFileInfo(rl), nil
}
