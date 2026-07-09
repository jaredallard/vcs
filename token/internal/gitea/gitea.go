// Copyright (C) 2026 vcs contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

// Package gitea implements [shared.Provider] for Gitea.
package gitea

import (
	"go.rgst.io/jaredallard/vcs/v2/token/internal/shared"
)

// Providers is a list of providers that can be used to retrieve a
// token for Gitea.
var Providers = []shared.Provider{
	&shared.EnvProvider{EnvVars: []shared.EnvVar{{Name: "GITEA_TOKEN"}}},
}
