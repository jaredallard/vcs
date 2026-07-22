package shared_test

import (
	"context"
	"testing"
	"time"

	"go.rgst.io/jaredallard/cmdexec/v2"
	"go.rgst.io/jaredallard/vcs/v3"
	"go.rgst.io/jaredallard/vcs/v3/token"
	"gotest.tools/v3/assert"
)

var bfalse = false

// clearHostToken clears the token for the host when fetching a Github
// token.
func clearHostToken(t *testing.T, newValue string) {
	cmdexec.UseMockExecutor(t, cmdexec.NewMockExecutor(&cmdexec.MockCommand{
		Name:   "gh",
		Args:   []string{"auth", "token"},
		Stdout: []byte("\n"),
	}))
	t.Setenv("GITHUB_TOKEN", newValue)
}

// TestCloneClonesAllAttributes ensures that Clone returns a new token
// with the same attributes as the original token.
func TestCloneClonesAllAttributes(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", time.Now().String())

	originalToken, err := token.Fetch(context.Background(), vcs.ProviderGithub, &token.Options{
		UseGlobalCache: &bfalse,
	})
	assert.NilError(t, err)
	assert.Assert(t, originalToken != nil, "expected a token to be returned")

	clone := originalToken.Clone()
	assert.DeepEqual(t, originalToken, clone)
}

func TestStringRedacts(t *testing.T) {
	clearHostToken(t, "token-xyz")

	originalToken, err := token.Fetch(context.Background(), vcs.ProviderGithub, &token.Options{
		UseGlobalCache: &bfalse,
	})
	assert.NilError(t, err)
	assert.Assert(t, originalToken != nil, "expected a token to be returned")

	assert.Equal(t, originalToken.String(), "toke*****", "expected token to be partially redacted")
}

func TestIsUnauthenticatedDetectsEmptyToken(t *testing.T) {
	clearHostToken(t, "")

	originalToken, err := token.Fetch(context.Background(), vcs.ProviderGithub, &token.Options{
		AllowUnauthenticated: true,
		UseGlobalCache:       &bfalse,
	})
	assert.NilError(t, err)
	assert.Assert(t, originalToken != nil, "expected a token to be returned")

	assert.Assert(t, originalToken.IsUnauthenticated(), "expected token to be unauthenticated")
}
