package token_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jaredallard/vcs"
	"github.com/jaredallard/vcs/token"
	"gotest.tools/v3/assert"
)

// ignoreTime is a [cmp.Option] that ignores time.Time values when
// comparing them, always returning true.
var ignoreTime = cmp.Comparer(func(_, _ time.Time) bool {
	// Times are random, so ignore them.
	return true
})

// TestCanGetToken ensures that [token.Fetch] calls the underlying
// provider to get the token.
func TestCanGetToken(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", time.Now().String())
	authToken, err := token.Fetch(context.Background(), vcs.ProviderGithub, false)
	assert.NilError(t, err)
	assert.Assert(t, authToken != nil, "expected a token to be returned")
	assert.DeepEqual(t, authToken, &token.Token{
		Source: "environment variable (GITHUB_TOKEN)",
		Value:  os.Getenv("GITHUB_TOKEN"),
	}, ignoreTime)
}

// TestCanGetCachedToken ensures that [token.Fetch] returns the same
// token when called multiple times and caching is enabled.
func TestCanGetCachedToken(t *testing.T) {
	bfalse := false
	t.Setenv("GITHUB_TOKEN", time.Now().String())

	originalToken, err := token.Fetch(context.Background(), vcs.ProviderGithub, false, &token.Options{UseGlobalCache: &bfalse})
	assert.NilError(t, err)
	assert.Assert(t, originalToken != nil, "expected a token to be returned")
	assert.DeepEqual(t, originalToken, &token.Token{
		Source: "environment variable (GITHUB_TOKEN)",
		Value:  os.Getenv("GITHUB_TOKEN"),
	}, ignoreTime)
	assert.Equal(t, originalToken.FetchedAt.IsZero(), false) // should not be zero

	// Fetch again, should return the same token.
	newToken, err := token.Fetch(context.Background(), vcs.ProviderGithub, false)
	assert.NilError(t, err)
	assert.Assert(t, newToken != nil, "expected a token to be returned")
	assert.DeepEqual(t, newToken, &token.Token{
		FetchedAt: originalToken.FetchedAt,
		Source:    "environment variable (GITHUB_TOKEN)",
		Value:     os.Getenv("GITHUB_TOKEN"),
	})
}
