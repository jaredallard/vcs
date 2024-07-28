package gitlab

import (
	"testing"

	"github.com/jaredallard/vcs/token/internal/shared"
	"gotest.tools/v3/assert"
)

// TestCanGetJobTokenFromEnv ensures that a job token can be read from
// the environment and that it has a type of TokenTypeJob.
func TestCanGetJobTokenFromEnv(t *testing.T) {
	t.Setenv("CI_JOB_TOKEN", "im-a-token")

	p := envProvider()

	got, err := p.Token()
	assert.NilError(t, err, "expected no error")
	assert.DeepEqual(t, &shared.Token{Value: "im-a-token", Type: TokenTypeJob}, got)
}
