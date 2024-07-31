// Copyright (C) 2024 vcs contributors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this program. If not, see
// <https://www.gnu.org/licenses/>.
//
// SPDX-License-Identifier: LGPL-3.0

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
