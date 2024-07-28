package shared_test

import (
	"testing"

	"github.com/jaredallard/vcs/token/internal/shared"
	"gotest.tools/v3/assert"
)

func TestEnvProviderReadsCorrectEnvVar(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "token")

	p := &shared.EnvProvider{EnvVars: []shared.EnvVar{{Name: "GITHUB_TOKEN"}}}
	token, err := p.Token()
	assert.NilError(t, err)
	assert.DeepEqual(t, &shared.Token{Value: "token"}, token)
}
