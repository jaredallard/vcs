// Copyright (C) 2026 vcs contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

// Description: Contains github specific Git functionality.

package git

import (
	"os"
	"path/filepath"
	"testing"

	"go.rgst.io/jaredallard/vcs/v2"
	"gotest.tools/v3/assert"
)

func Test_cloneArchive(t *testing.T) {
	type args struct {
		vcsp      vcs.Provider
		ref       string
		sourceURL string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "can download archive",
			args: args{
				ref:       "main",
				sourceURL: "https://github.com/jaredallard/vcs",
			},
		},
		{
			name: "supports git urls",
			args: args{
				ref:       "v0.1.0",
				sourceURL: "git://github.com/jaredallard/vcs",
			},
		},
		{
			name: "supports ssh urls",
			args: args{
				ref:       "v0.2.0",
				sourceURL: "git@github.com:jaredallard/vcs",
			},
		},
		{
			name: "supports .git at end of url",
			args: args{
				ref:       "v0.2.0",
				sourceURL: "https://github.com/jaredallard/vcs.git",
			},
		},
		{
			name: "can download gitea (forgejo) archive",
			args: args{
				ref:       "main",
				sourceURL: "https://git.rgst.io/jaredallard/vcs",
				vcsp:      vcs.ProviderForgejo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.vcsp == "" {
				tt.args.vcsp = vcs.ProviderGithub
			}

			got, err := cloneArchive(t.Context(), tt.args.vcsp, tt.args.ref, tt.args.sourceURL, t.TempDir())
			if (err != nil) != tt.wantErr {
				t.Errorf("cloneArchiveGithub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// debugging information for the next check
			files, err := os.ReadDir(got)
			assert.NilError(t, err, "failed to read directory")
			if len(files) > 0 {
				t.Logf("Directory contains %d file(s)/directories:", len(files))
				t.Log("=================================================")
			}
			for _, f := range files {
				t.Logf("%s", f.Name())
			}
			if len(files) > 0 {
				t.Log("=================================================")
			}

			// Ensure that there's a file in the directory.
			_, err = os.Stat(filepath.Join(got, "README.md"))
			assert.NilError(t, err, "expected README.md to exist in the archive")
		})
	}
}
