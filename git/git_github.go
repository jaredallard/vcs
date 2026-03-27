// Copyright (C) 2026 vcs contributors
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
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"code.gitea.io/sdk/gitea"
	giturls "github.com/chainguard-dev/git-urls"
	"github.com/google/go-github/v82/github"
	"go.rgst.io/jaredallard/archives/v2"
	"go.rgst.io/jaredallard/vcs/v2"
	"go.rgst.io/jaredallard/vcs/v2/token"
)

// cloneArchive is the same as [Clone] but uses the Provider API to
// download the repository contents at a specific ref. These archives do
// not contain the .git directory and thus may not always be desirable.
func cloneArchive(ctx context.Context, vcsp vcs.Provider, ref, sourceURL, tempDir string) (string, error) {
	u, err := giturls.Parse(sourceURL)
	if err != nil {
		return "", err
	}

	t, err := token.Fetch(ctx, vcsp, false, &token.Options{
		AllowUnauthenticated: true,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get github token for archive fetch: %w", err)
	}

	owner, repo := path.Split(u.Path)

	// Attempt to normalize the owner and repo just in case.
	owner = strings.ReplaceAll(owner, "/", "")
	repo = strings.TrimSuffix(repo, ".git")

	var body io.ReadCloser
	switch vcsp { //nolint:exhaustive // Why: see default
	case vcs.ProviderGithub:
		gh := github.NewClient(nil)
		if !t.IsUnauthenticated() {
			gh = gh.WithAuthToken(t.Value)
		}

		rc, _, err := gh.Repositories.GetArchiveLink(ctx, owner, repo, github.Tarball, &github.RepositoryContentGetOptions{
			Ref: ref,
		}, 0)
		if err != nil {
			return "", fmt.Errorf("failed to get archive link: %w", err)
		}

		req, err := gh.NewRequest(http.MethodGet, rc.String(), http.NoBody)
		if err != nil {
			return "", fmt.Errorf("failed to download archive: %w", err)
		}

		resp, err := gh.BareDo(ctx, req)
		if err != nil {
			return "", fmt.Errorf("failed to download archive: %w", err)
		}
		body = resp.Body
	case vcs.ProviderForgejo, vcs.ProviderGitea:
		gsdk, err := gitea.NewClient(u.Scheme + "://" + u.Host)
		if err != nil {
			return "", fmt.Errorf("failed to create gitea sdk: %w", err)
		}
		if !t.IsUnauthenticated() {
			gsdk.SetBasicAuth("token", t.Value)
		}

		body, _, err = gsdk.GetArchiveReader(owner, repo, ref, gitea.TarGZArchive)
		if err != nil {
			return "", fmt.Errorf("failed to download archive: %w", err)
		}
	default:
		return "", fmt.Errorf("provider %q doesn't support archive downloads", vcsp)
	}
	defer body.Close() //nolint:errcheck,gosec // Why: Best effort.

	if err := archives.Extract(body, tempDir, archives.ExtractOptions{Extension: ".tar.gz"}); err != nil {
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

		// Should either contain the org+repo, OR be exactly equal to the
		// repo name (forgejo/gitea)
		hasOwnerAndRepo := strings.Contains(f.Name(), owner) && strings.Contains(f.Name(), repo)
		exactRepo := f.Name() == repo
		if !hasOwnerAndRepo && !exactRepo {
			continue
		}

		dir = f.Name()
		break
	}

	return filepath.Join(tempDir, dir), nil
}
