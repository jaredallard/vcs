// Copyright (C) 2026 vcs contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestRedactURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "redacts http user and password",
			url:  "https://user:pass@github.com/org/repo",
			want: "https://redacted@github.com/org/repo",
		},
		{
			name: "redacts http token as username",
			url:  "https://s3cr3t@github.com/org/repo",
			want: "https://redacted@github.com/org/repo",
		},
		{
			name: "leaves urls without credentials unchanged",
			url:  "https://github.com/org/repo",
			want: "https://github.com/org/repo",
		},
		{
			name: "leaves scp-style urls unchanged",
			url:  "git@github.com:org/repo",
			want: "git@github.com:org/repo",
		},
		{
			name: "leaves ssh username unchanged",
			url:  "ssh://git@github.com/org/repo",
			want: "ssh://git@github.com/org/repo",
		},
		{
			name: "redacts ssh password",
			url:  "ssh://git:pass@github.com/org/repo",
			want: "ssh://git:redacted@github.com/org/repo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, redactURL(tt.url), tt.want)
		})
	}
}

func TestRedactErr(t *testing.T) {
	origErr := &exec.ExitError{}

	err := redactErr(fmt.Errorf("failed to clone https://user:s3cr3t@host/repo: %w", origErr), "https://user:s3cr3t@host/repo")
	assert.Equal(t, err.Error(), "failed to clone https://redacted@host/repo: "+origErr.Error())

	var execErr *exec.ExitError
	assert.Assert(t, errors.As(err, &execErr), "expected original error to be accessible")

	noCredsErr := errors.New("failed to clone https://host/repo")
	assert.Equal(t, redactErr(noCredsErr, "https://host/repo"), noCredsErr)
	assert.NilError(t, redactErr(nil, "https://user:s3cr3t@host/repo"))
}

func TestCloneDoesNotLeakCredentials(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("TMPDIR", tmpDir)

	_, err := Clone(t.Context(), "main", "https://user:s3cr3t@nonexistent.invalid/org/repo")
	assert.Assert(t, err != nil, "expected clone to fail")
	assert.Assert(t, !strings.Contains(err.Error(), "s3cr3t"), "error contains credentials: %v", err)

	entries, err := os.ReadDir(tmpDir)
	assert.NilError(t, err)
	for _, e := range entries {
		assert.Assert(t, !strings.Contains(e.Name(), "s3cr3t"), "temp dir contains credentials: %s", e.Name())
	}
}

func TestCloneRemovesTempDirOnFailure(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("TMPDIR", tmpDir)

	_, err := Clone(t.Context(), "main", "https://nonexistent.invalid/org/repo")
	assert.Assert(t, err != nil, "expected clone to fail")

	entries, err := os.ReadDir(tmpDir)
	assert.NilError(t, err)
	assert.Equal(t, len(entries), 0, "expected temp dir to be removed, found %v", entries)
}
