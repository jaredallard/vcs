// Copyright (C) 2026 vcs contributors
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
