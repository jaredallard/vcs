package github_test

import (
	"testing"

	"github.com/jaredallard/vcs/token/github"
	"gotest.tools/v3/assert"
)

func TestEnvProviderReadsCorrectEnvVar(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "token")

	p := &github.EnvProvider{}
	token, err := p.Token()
	assert.NilError(t, err)
	assert.Equal(t, "token", token)
}
