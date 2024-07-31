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

// Package opts contains the options and interfaces for the releases
// package. Stored separately to avoid circular dependencies.
package opts

import (
	"context"
	"io"
	"os"

	"github.com/jaredallard/vcs"
	"github.com/jaredallard/vcs/token"
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
