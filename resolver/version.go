// Copyright (C) 2024 Jared Allard <jaredallard@users.noreply.github.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package resolver

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

// Version represents a version found in a Git repository. Versions are
// only discovered if a tag or branch points to a commit (individual
// commits will never be automatically discovered unless they are
// manually passed in).
type Version struct {
	// Commit is the underlying commit hash for this version.
	Commit string `yaml:"commit,omitempty"`

	// Tag is the underlying tag for this version, if set.
	Tag string `yaml:"tag,omitempty"`
	sv  *semver.Version

	// Virtual denotes that a version does not actually relate to anything
	// in a Git repository. This is mainly meant for testing or scenarios
	// where you want to inject something in that is not a real version.
	//
	// This field is never set by the resolver in this package.
	Virtual string `yaml:"virtual,omitempty"`

	// Branch is the underlying branch for this version, if set.
	Branch string `yaml:"branch,omitempty"`
}

// Equal returns true if the two versions are equal.
func (v *Version) Equal(other *Version) bool {
	// If either is nil, they must both be nil.
	if v == nil || other == nil {
		return v == other
	}

	// Otherwise, check all fields.
	return v.Commit == other.Commit && v.Tag == other.Tag && v.Branch == other.Branch
}

// String is a user-friendly representation of the version that can be
// used in error messages.
func (v *Version) String() string {
	switch {
	case v.Virtual != "":
		return fmt.Sprintf("virtual (source: %s)", v.Virtual)
	case v.Tag != "":
		return fmt.Sprintf("tag %s (%s)", v.Tag, v.Commit)
	case v.Branch != "":
		return fmt.Sprintf("branch %s (%s)", v.Branch, v.Commit)
	default:
		return v.Commit
	}
}

// GitRef returns a Git reference that can be used to check out the
// version.
func (v *Version) GitRef() string {
	switch {
	// TODO(jaredallard): This will require native ext handling.
	case v.Virtual != "":
		return "NOT_A_VALID_GIT_VERSION"
	case v.Tag != "":
		return "refs/tags/" + v.Tag
	case v.Branch != "":
		return "refs/heads/" + v.Branch
	default:
		return v.Commit
	}
}
