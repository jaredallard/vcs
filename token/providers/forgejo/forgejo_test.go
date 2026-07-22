package forgejo_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.rgst.io/jaredallard/cmdexec/v2"
	"go.rgst.io/jaredallard/vcs/v3/token"
	"go.rgst.io/jaredallard/vcs/v3/token/providers/forgejo"
	"gotest.tools/v3/assert"
)

// TestCanGetToken ensures that the provider can get a token.
func TestCanGetToken(t *testing.T) {
	p := &forgejo.FJProvider{}

	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	fakeValue := time.Now().Format(time.UnixDate)

	b, err := json.Marshal(forgejo.FJKeys{
		Hosts: map[string]forgejo.FJKeysHost{
			"git.rgst.io": {Token: fakeValue},
		},
	})
	assert.NilError(t, err, "expected marshal keys.json to not fail")

	keysPath := filepath.Join(tmpDir, ".local", "share", "forgejo-cli", "keys.json")
	assert.NilError(t,
		os.MkdirAll(filepath.Dir(keysPath), 0o755),
		"expected dir %s create not fail", filepath.Dir(keysPath),
	)

	assert.NilError(t, os.WriteFile(keysPath, b, 0o655), "expected write %s not fail", keysPath)

	cmdexec.UseMockExecutor(t, cmdexec.NewMockExecutor(&cmdexec.MockCommand{
		Name:   "fj",
		Args:   []string{"whoami"},
		Stdout: []byte("currently signed in to not-a-real-user@git.rgst.io"),
	}))

	got, err := p.Token()
	assert.NilError(t, err)
	assert.DeepEqual(t, &token.Token{
		Source: "fj",
		Value:  fakeValue,
	}, got)
}
