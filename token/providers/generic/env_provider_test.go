package generic_test

import (
	"fmt"
	"testing"

	"go.rgst.io/jaredallard/vcs/v3/token/internal/shared"
	"go.rgst.io/jaredallard/vcs/v3/token/providers/generic"
	"gotest.tools/v3/assert"
)

func TestEnvProviderReadsCorrectEnvVar(t *testing.T) {
	t.Setenv(t.Name(), "token")

	p := &generic.EnvProvider{EnvVars: []generic.EnvVar{{Name: t.Name()}}}
	tok, err := p.Token()
	assert.NilError(t, err)
	assert.DeepEqual(t, &shared.Token{
		Source: fmt.Sprintf("environment variable (%s)", t.Name()),
		Value:  "token",
	}, tok)
}
