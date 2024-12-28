package gitlab

import (
	"testing"

	"github.com/jaredallard/cmdexec"
	"github.com/jaredallard/vcs/token/internal/shared"
	"gotest.tools/v3/assert"
)

// TestTrimsSpace ensures that the token returned by the glabProvider is
// trimmed of any leading or trailing whitespace.
func TestTrimsSpace(t *testing.T) {
	p := &GlabProvider{}

	cmdexec.UseMockExecutor(t, cmdexec.NewMockExecutor(
		&cmdexec.MockCommand{
			Name:   "glab",
			Args:   []string{"config", "get", "-g", "host"},
			Stdout: []byte("gitlab.com\n"),
		},
		&cmdexec.MockCommand{
			Name:   "glab",
			Args:   []string{"config", "get", "-g", "token", "-h", "gitlab.com"},
			Stdout: []byte(" token\n"),
		},
	))

	got, err := p.Token()
	assert.NilError(t, err)
	assert.DeepEqual(t, &shared.Token{
		Source: "glab",
		Value:  "token",
	}, got)
}

// TestCanGetJobTokenFromEnv ensures that a job token can be read from
// the environment and that it has a type of TokenTypeJob.
func TestCanGetJobTokenFromEnv(t *testing.T) {
	t.Setenv("CI_JOB_TOKEN", "im-a-token")

	p := envProvider()

	got, err := p.Token()
	assert.NilError(t, err, "expected no error")
	assert.DeepEqual(t, &shared.Token{
		Source: "environment variable (CI_JOB_TOKEN)",
		Value:  "im-a-token",
		Type:   TokenTypeJob,
	}, got)
}
