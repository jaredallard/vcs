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

// Package git contains functions for interacting with Git repositories
// using the Git CLI. As such, this package requires the Git CLI to be
// installed on the system.
package git

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/jaredallard/cmdexec"
	"github.com/jaredallard/vcs"
	"github.com/jaredallard/vcs/internal/execerr"
	"github.com/pkg/errors"
)

// This block contains errors and regexes
var (
	// ErrNoHeadBranch is returned when a repository's HEAD (aka default) branch cannot
	// be determine
	ErrNoHeadBranch = errors.New("failed to find a head branch, does one exist?")

	// ErrNoRemoteHeadBranch is returned when a repository's remote  default/HEAD branch
	// cannot be determined.
	ErrNoRemoteHeadBranch = errors.New("failed to get head branch from remote origin")

	// headPattern is used to parse git output to determine the head branch
	headPattern = regexp.MustCompile(`HEAD branch: ([[:alpha:]]+)`)
)

// gitCommand is a wrapper for [cmdexec.CommandContext] for git which
// forces `LC_ALL=C` into the environment, in the event that the
// command output is parsed.
func gitCommand(ctx context.Context, args ...string) cmdexec.Cmd {
	cmd := cmdexec.CommandContext(ctx, "git", args...)
	env := os.Environ()
	env = append(env, "LC_ALL=C")
	cmd.SetEnviron(env)
	return cmd
}

// GetDefaultBranch determines the default/HEAD branch for a given git
// repository.
func GetDefaultBranch(ctx context.Context, path string) (string, error) {
	cmd := gitCommand(ctx, "remote", "show", "origin")
	cmd.SetDir(path)
	out, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "failed to get head branch from remote origin")
	}

	matches := headPattern.FindStringSubmatch(string(out))
	if len(matches) != 2 {
		return "", ErrNoRemoteHeadBranch
	}

	return matches[1], nil
}

// CloneOptions contains options accepted by [Clone].
type CloneOptions struct {
	// UseArchive fetches the references using a tarball from the
	// underlying VCS provider. This comes with the caveat that the
	// downloaded tarball will not contain the .git directory. However,
	// this can be much faster to download than a full clone.
	//
	// Currently, only Github URLs are supported.
	//
	// If this option fails, a normal clone will be performed without an
	// error.
	UseArchive bool
}

// Clone clone a git repository to a temporary directory and returns the
// path to the repository. If ref is empty, the default branch will be
// used. A shallow clone is performed.
//
// optss is a variadic argument only to avoid a breaking change. Only
// one option struct is allowed, an error will be returned if more than
// one is provided.
func Clone(ctx context.Context, ref, url string, optss ...*CloneOptions) (string, error) {
	tempDir, err := os.MkdirTemp("", strings.ReplaceAll(url, "/", "-"))
	if err != nil {
		return "", errors.Wrap(err, "failed to create temporary directory")
	}

	// Read opts from the variadic argument. We use a variadic argument
	// here to avoid a breaking change when this variable was added.
	var opts CloneOptions
	if len(optss) == 1 {
		if optss[0] != nil {
			opts = *optss[0]
		}
	} else if len(optss) > 1 {
		return "", fmt.Errorf("too many options provided")
	}

	if opts.UseArchive {
		provider, err := vcs.ProviderFromURL(url, nil)
		if err == nil && provider == vcs.ProviderGithub {
			tmpDir, err := cloneArchiveGithub(ctx, ref, url, tempDir)
			if err == nil {
				return tmpDir, nil
			}
		}
	}

	cmds := [][]string{
		{"init"},
		{"remote", "add", "origin", url},
		{"-c", "protocol.version=2", "fetch", "origin", ref},
		{"reset", "--hard", "FETCH_HEAD"},
	}
	for _, cmd := range cmds {
		//nolint:gosec // Why: Commands are not user provided.
		c := gitCommand(ctx, cmd...)
		c.SetDir(tempDir)
		if err := c.Run(); err != nil {
			var execErr *exec.ExitError
			if errors.As(err, &execErr) {
				return "", fmt.Errorf("failed to run %q (%w): %s", cmd, err, string(execErr.Stderr))
			}

			return "", fmt.Errorf("failed to run %q: %w", cmd, err)
		}
	}

	return tempDir, nil
}

// ListRemote returns a list of all remotes as shown from running 'git
// ls-remote'.
func ListRemote(ctx context.Context, remote string) ([][]string, error) {
	cmd := gitCommand(ctx, "ls-remote", remote)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get remote branches: %w", execerr.From(err))
	}

	remotes := make([][]string, 0)

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		remotes = append(remotes, strings.Fields(line))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return remotes, nil
}
