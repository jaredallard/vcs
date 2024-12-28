package github_test

import (
	"testing"

	"github.com/jaredallard/cmdexec"
	"github.com/jaredallard/vcs/token"
	"github.com/jaredallard/vcs/token/internal/github"
	"gotest.tools/v3/assert"
)

// TestTrimsSpace ensures that the token returned by the ghProvider is
// trimmed of any leading or trailing whitespace.
func TestTrimsSpace(t *testing.T) {
	p := &github.GHProvider{}

	cmdexec.UseMockExecutor(t, cmdexec.NewMockExecutor(&cmdexec.MockCommand{
		Name:   "gh",
		Args:   []string{"auth", "token"},
		Stdout: []byte(" token\n"),
	}))

	got, err := p.Token()
	assert.NilError(t, err)
	assert.DeepEqual(t, &token.Token{
		Source: "gh",
		Value:  "token",
	}, got)
}
