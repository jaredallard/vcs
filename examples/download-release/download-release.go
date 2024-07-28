package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/jaredallard/vcs/releases"
	"github.com/jaredallard/vcs/resolver"
)

// main downloads a release of stencil from Github to the local
// directory. It uses the GOOS and GOARCH of the current system to
// determine which asset to download.
//
// Under the hood, the VCS provider is being determined from the URL and
// authentication is being provided for the determined VCS provider if
// it was configured on the system.
func main() {
	ctx := context.Background()

	repoURL := "https://github.com/rgst-io/stencil"

	r := resolver.NewResolver()
	v, err := r.Resolve(ctx, repoURL, &resolver.Criteria{Constraint: "*"})
	if err != nil {
		panic(err)
	}

	fmt.Println("Latest version of stencil is", v)

	fmt.Println("Downloading release...")
	resp, fi, err := releases.Fetch(ctx, &releases.FetchOptions{
		RepoURL:   "https://github.com/rgst-io/stencil",
		AssetName: fmt.Sprintf("stencil_*_%s_%s.tar.gz", runtime.GOOS, runtime.GOARCH),
		Tag:       v.Tag,
	})
	if err != nil {
		panic(err)
	}
	defer resp.Close()

	f, err := os.Create(fi.Name())
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp); err != nil {
		panic(err)
	}

	fmt.Println("Downloaded release to", fi.Name())
}
