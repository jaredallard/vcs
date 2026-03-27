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
