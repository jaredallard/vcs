package releases

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"testing"

	"github.com/jaredallard/vcs/releases/internal/opts"
)

func TestFetch(t *testing.T) {
	type args struct {
		opts *opts.FetchOptions
	}
	tests := []struct {
		name string
		args args
		// want is a hash of the expected output
		want     string
		wantName string
		wantErr  bool
	}{
		{
			name: "should fetch a github release",
			args: args{
				opts: &opts.FetchOptions{
					RepoURL:   "https://github.com/rgst-io/stencil",
					Tag:       "v0.7.0",
					AssetName: "stencil_0.7.0_linux_arm64.tar.gz",
				},
			},
			want:     "68992b329703c8579fc063932975f5aae45157a4c5c19eb0364c3b153e08a106",
			wantName: "stencil_0.7.0_linux_arm64.tar.gz",
			wantErr:  false,
		},
		{
			name: "should fetch the correct asset when given a list",
			args: args{
				opts: &opts.FetchOptions{
					RepoURL:    "https://github.com/rgst-io/stencil",
					Tag:        "v0.7.0",
					AssetNames: []string{"stencil_0.7.0_linux_arm64.tar.gz"},
				},
			},
			want:     "68992b329703c8579fc063932975f5aae45157a4c5c19eb0364c3b153e08a106",
			wantName: "stencil_0.7.0_linux_arm64.tar.gz",
			wantErr:  false,
		},
		{
			name: "should fail when given an invalid tag",
			args: args{
				opts: &opts.FetchOptions{
					RepoURL: "https://github.com/rgst-io/stencil",
					Tag:     "i-am-not-a-real-tag",
				},
			},
			wantErr: true,
		},
		{
			name: "should fail when given an invalid repo URL",
			args: args{
				opts: &opts.FetchOptions{
					RepoURL: "not-a-real-repo-url",
					Tag:     "a-tag",
				},
			},
			wantErr: true,
		},
		{
			name: "should fail when no asset given",
			args: args{
				opts: &opts.FetchOptions{
					RepoURL: "https://github.com/rgst-io/stencil",
					Tag:     "v0.7.0",
				},
			},
			wantErr: true,
		},
		{
			name: "should fail when no repo URL given",
			args: args{
				opts: &opts.FetchOptions{},
			},
			wantErr: true,
		},
		{
			name: "should fail when no tag given",
			args: args{
				opts: &opts.FetchOptions{
					RepoURL: "a-repo",
				},
			},
			wantErr: true,
		},
		{
			name:    "should fail when no opts given",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, fi, err := Fetch(context.Background(), "", tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			defer got.Close()

			b, err := io.ReadAll(got)
			if err != nil {
				t.Errorf("Fetch() error = %v", err)
				return
			}

			hashByt := sha256.Sum256(b)
			hash := hex.EncodeToString(hashByt[:])
			if hash != tt.want {
				t.Errorf("Fetch() hash = %v, want %v", hash, tt.want)
			}
			if fi.Name() != tt.wantName {
				t.Errorf("Fetch() name = %v, wantName %v", fi.Name(), tt.wantName)
			}
		})
	}
}

func TestGetReleaseNotes(t *testing.T) {
	type args struct {
		opts *opts.GetReleaseNoteOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should get release notes for a github release",
			args: args{
				opts: &opts.GetReleaseNoteOptions{
					RepoURL: "https://github.com/rgst-io/stencil",
					Tag:     "v0.7.0",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetReleaseNotes(context.Background(), "", tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReleaseNotes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if got == "" {
				t.Errorf("GetReleaseNotes() return empty string")
			}
		})
	}
}
