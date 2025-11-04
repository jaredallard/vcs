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

// Description: Contains github specific Git functionality.

package git

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	giturls "github.com/chainguard-dev/git-urls"
	"github.com/google/go-github/v77/github"
	"github.com/jaredallard/archives"
	"github.com/jaredallard/vcs"
	"github.com/jaredallard/vcs/token"
)

// cloneArchiveGithub is the same as [Clone] but uses the Github API to
// download the repository contents at a specific ref. These archives do
// not contain the .git directory and thus may not always be desirable.
func cloneArchiveGithub(ctx context.Context, ref, sourceURL, tempDir string) (string, error) {
	u, err := giturls.Parse(sourceURL)
	if err != nil {
		return "", err
	}

	t, err := token.Fetch(ctx, vcs.ProviderGithub, true)
	if err != nil {
		return "", fmt.Errorf("failed to get github token for archive fetch: %w", err)
	}

	gh := github.NewClient(nil).WithAuthToken(t.Value)

	owner, repo := path.Split(u.Path)

	// Attempt to normalize the owner and repo just in case.
	owner = strings.ReplaceAll(owner, "/", "")
	repo = strings.TrimSuffix(repo, ".git")

	rc, _, err := gh.Repositories.GetArchiveLink(ctx, owner, repo, github.Tarball, &github.RepositoryContentGetOptions{
		Ref: ref,
	}, 0)
	if err != nil {
		return "", fmt.Errorf("failed to get archive link: %w", err)
	}

	resp, err := http.Get(rc.String())
	if err != nil {
		return "", fmt.Errorf("failed to download archive: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck,gosec // Why: Best effort.

	if err := archives.Extract(resp.Body, tempDir, archives.ExtractOptions{Extension: ".tar.gz"}); err != nil {
		return "", fmt.Errorf("failed to extract archive: %w", err)
	}

	// The extracted archive contains a top-level directory in it, so
	// select the first directory in the tempDir.
	files, err := os.ReadDir(tempDir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	// Select the first directory in the tempDir.
	var dir string
	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		// Should contain the owner and repo name in it.
		//nolint:staticcheck // Why: This is easy enough to read.
		if !(strings.Contains(f.Name(), owner) && strings.Contains(f.Name(), repo)) {
			continue
		}

		dir = f.Name()
		break
	}

	return filepath.Join(tempDir, dir), nil
}
