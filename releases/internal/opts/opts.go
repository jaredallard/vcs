// Copyright (C) 2026 vcs contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

// Package opts contains the options and interfaces for the releases
// package. Stored separately to avoid circular dependencies.
package opts

import (
	"context"
	"io"
	"os"

	"go.rgst.io/jaredallard/vcs/v3"
	"go.rgst.io/jaredallard/vcs/v3/token"
)

// Fetcher is an interface that fetches assets from a release. VCS
// providers must implement this interface.
type Fetcher interface {
	// Fetch returns an asset as a io.ReadCloser
	Fetch(ctx context.Context, token *token.Token, opts *FetchOptions) (io.ReadCloser, os.FileInfo, error)

	// GetReleaseNotes returns the release notes of a release
	GetReleaseNotes(ctx context.Context, token *token.Token, opts *GetReleaseNoteOptions) (string, error)
}

// FetchOptions is a set of options for Fetch
type FetchOptions struct {
	Overrides []vcs.Override

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
	Overrides []vcs.Override

	// RepoURL is the repository URL, it should be a valid
	// URL.
	RepoURL string

	// Tag is the tag of the release
	Tag string
}
