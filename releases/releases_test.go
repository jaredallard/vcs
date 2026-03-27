package releases

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"testing"

	"go.rgst.io/jaredallard/vcs/v2"
)

func TestFetch(t *testing.T) {
	type args struct {
		opts *FetchOptions
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
				opts: &FetchOptions{
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
				opts: &FetchOptions{
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
				opts: &FetchOptions{
					RepoURL: "https://github.com/rgst-io/stencil",
					Tag:     "i-am-not-a-real-tag",
				},
			},
			wantErr: true,
		},
		{
			name: "should fail when given an invalid repo URL",
			args: args{
				opts: &FetchOptions{
					RepoURL: "not-a-real-repo-url",
					Tag:     "a-tag",
				},
			},
			wantErr: true,
		},
		{
			name: "should fail when no asset given",
			args: args{
				opts: &FetchOptions{
					RepoURL: "https://github.com/rgst-io/stencil",
					Tag:     "v0.7.0",
				},
			},
			wantErr: true,
		},
		{
			name: "should fail when no repo URL given",
			args: args{
				opts: &FetchOptions{},
			},
			wantErr: true,
		},
		{
			name: "should fail when no tag given",
			args: args{
				opts: &FetchOptions{
					RepoURL: "a-repo",
				},
			},
			wantErr: true,
		},
		{
			name:    "should fail when no opts given",
			wantErr: true,
		},
		{
			name: "should support gitlab",
			args: args{
				opts: &FetchOptions{
					RepoURL:   "https://gitlab.com/jaredallard/vcs-test-repo",
					Tag:       "v0.1.0",
					AssetName: "vcs-test-repo-v0.1.0.tar.gz",
				},
			},
			want:     "ab68a8000fde646c0331118e0e54c0257fa86bc1a9e308f09153d16024a8eb61",
			wantName: "vcs-test-repo-v0.1.0.tar.gz",
		},
		{
			name: "should support gitea (forgejo)",
			args: args{
				opts: &FetchOptions{
					RepoURL:   "https://git.rgst.io/jaredallard/ingress-anubis",
					Tag:       "v1.10.4",
					AssetName: "checksums.txt",
				},
			},
			want:     "df62dbe653de600d5d76ae2dcb03baf10ea6081b9b6cf8347e0842f81e21e16c",
			wantName: "checksums.txt",
		},
		{
			name: "should support overrides",
			args: args{
				opts: &FetchOptions{
					Overrides: []vcs.Override{
						// Set Gitlab to the wrong provider to force it to fail.
						{URLBase: "https://gitlab.com", Provider: vcs.ProviderGithub},
					},
					RepoURL:   "https://gitlab.com/jaredallard/vcs-test-repo",
					Tag:       "v0.1.0",
					AssetName: "vcs-test-repo-v0.1.0.tar.gz",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, fi, err := Fetch(t.Context(), tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			defer got.Close()

			h := sha256.New()
			if _, err := io.Copy(h, got); err != nil {
				t.Errorf("Fetch() error = %v", err)
				return
			}

			hash := hex.EncodeToString(h.Sum(nil))
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
		opts *GetReleaseNoteOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should get release notes for a github release",
			args: args{
				opts: &GetReleaseNoteOptions{
					RepoURL: "https://github.com/rgst-io/stencil",
					Tag:     "v0.7.0",
				},
			},
			wantErr: false,
		},
		{
			name: "should get release notes for a gitlab release",
			args: args{
				opts: &GetReleaseNoteOptions{
					RepoURL: "https://gitlab.com/jaredallard/vcs-test-repo",
					Tag:     "v0.1.0",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := GetReleaseNotes(context.Background(), tt.args.opts)
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
