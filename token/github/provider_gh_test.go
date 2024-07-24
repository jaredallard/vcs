package github_test

import (
	"testing"

	"github.com/jaredallard/cmdexec"
	"github.com/jaredallard/vcs/token/github"
	"gotest.tools/v3/assert"
)

// TestGhProviderTrimsSpace ensures that the token returned by the
// ghProvider is trimmed of any leading or trailing whitespace.
func TestGhProviderTrimsSpace(t *testing.T) {
	p := &github.GHProvider{}

	cmdexec.UseMockExecutor(t, cmdexec.NewMockExecutor(&cmdexec.MockCommand{
		Name:   "gh",
		Args:   []string{"auth", "token"},
		Stdout: []byte(" token\n"),
	}))

	token, err := p.Token()
	assert.NilError(t, err)
	assert.Equal(t, "token", token)
}
