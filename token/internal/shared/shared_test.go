package shared_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jaredallard/cmdexec"
	"github.com/jaredallard/vcs"
	"github.com/jaredallard/vcs/token"
	"github.com/jaredallard/vcs/token/internal/shared"
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

func TestEnvProviderReadsCorrectEnvVar(t *testing.T) {
	t.Setenv(t.Name(), "token")

	p := &shared.EnvProvider{EnvVars: []shared.EnvVar{{Name: t.Name()}}}
	tok, err := p.Token()
	assert.NilError(t, err)
	assert.DeepEqual(t, &shared.Token{
		Source: fmt.Sprintf("environment variable (%s)", t.Name()),
		Value:  "token",
	}, tok)
}

// TestCloneClonesAllAttributes ensures that Clone returns a new token
// with the same attributes as the original token.
func TestCloneClonesAllAttributes(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", time.Now().String())

	originalToken, err := token.Fetch(context.Background(), vcs.ProviderGithub, false, &token.Options{
		UseGlobalCache: &bfalse,
	})
	assert.NilError(t, err)
	assert.Assert(t, originalToken != nil, "expected a token to be returned")

	clone := originalToken.Clone()
	assert.DeepEqual(t, originalToken, clone)
}

func TestStringRedacts(t *testing.T) {
	clearHostToken(t, "token-xyz")

	originalToken, err := token.Fetch(context.Background(), vcs.ProviderGithub, false, &token.Options{
		UseGlobalCache: &bfalse,
	})
	assert.NilError(t, err)
	assert.Assert(t, originalToken != nil, "expected a token to be returned")

	assert.Equal(t, originalToken.String(), "toke*****", "expected token to be partially redacted")
}

func TestIsUnauthenticatedDetectsEmptyToken(t *testing.T) {
	clearHostToken(t, "")

	originalToken, err := token.Fetch(context.Background(), vcs.ProviderGithub, false, &token.Options{
		AllowUnauthenticated: true,
		UseGlobalCache:       &bfalse,
	})
	assert.NilError(t, err)
	assert.Assert(t, originalToken != nil, "expected a token to be returned")

	assert.Assert(t, originalToken.IsUnauthenticated(), "expected token to be unauthenticated")
}
