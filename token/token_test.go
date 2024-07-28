package token_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jaredallard/vcs"
	"github.com/jaredallard/vcs/token"
	"gotest.tools/v3/assert"
)

// TestCanGetToken ensures that [token.Fetch] calls the underlying
// provider to get the token.
func TestCanGetToken(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", time.Now().String())
	authToken, err := token.Fetch(context.Background(), vcs.ProviderGithub, false)
	assert.NilError(t, err)
	assert.Assert(t, authToken != nil, "expected a token to be returned")
	assert.DeepEqual(t, authToken, &token.Token{Value: os.Getenv("GITHUB_TOKEN")})
}
