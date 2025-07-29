//go:build !test_no_internet

package git_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jaredallard/vcs/git"
	"gotest.tools/v3/assert"
)

func TestGit(t *testing.T) {
	ctx := t.Context()

	// test clone
	dir, err := git.Clone(ctx, "main", "https://github.com/jaredallard/jaredallard")
	assert.NilError(t, err)

	assert.Assert(t, dir != "", "expected a directory to be returned")

	// ensure .git exists in the directory
	_, err = os.Stat(filepath.Join(dir, ".git"))
	assert.NilError(t, err, "expected .git to exist in the cloned directory")

	t.Run("GetDefaultBranch", func(t *testing.T) {
		t.Parallel()

		branch, err := git.GetDefaultBranch(ctx, dir)
		assert.NilError(t, err)
		assert.Equal(t, "main", branch)
	})

	t.Run("ListRemote", func(t *testing.T) {
		t.Parallel()

		remotes, err := git.ListRemote(ctx, dir)
		assert.NilError(t, err)
		assert.Assert(t, len(remotes) > 0, "expected at least one remote")
	})

	t.Run("Clone_Opts_UseArchive", func(t *testing.T) {
		t.Parallel()

		dir, err := git.Clone(ctx, "v0.2.0", "https://github.com/jaredallard/vcs", &git.CloneOptions{UseArchive: true})
		assert.NilError(t, err)

		// ensure .git does not exist in the directory
		_, err = os.Stat(filepath.Join(dir, ".git"))
		assert.ErrorContains(t, err, "no such file or directory")
	})
}

// Makes sure that GetDefaultBranch works correctly even when the system language is not set to English.
// Not parallel because it sets an environment variable.
func TestGetDefaultBranchDifferentOSLanguage(t *testing.T) {
	ctx := t.Context()

	dir, err := git.Clone(ctx, "main", "https://github.com/jaredallard/jaredallard")
	assert.NilError(t, err)

	assert.Assert(t, dir != "", "expected a directory to be returned")

	t.Setenv("LC_ALL", "fr_FR.UTF-8")
	defaultBranch, err := git.GetDefaultBranch(ctx, dir)
	assert.NilError(t, err)
	assert.Equal(t, defaultBranch, "main", "Expected default branch to be 'main'")
}
