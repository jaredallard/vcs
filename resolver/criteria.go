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
	"regexp"
	"strings"
	"sync"

	"github.com/Masterminds/semver/v3"
)

// constRexp is a regular expression that matches any non-numeric or "v"
// characters. Used to strip constraints to convert them into "versions".
var constRexp = regexp.MustCompile(`^[^v\d]+`)

// Criteria represents a set of criteria that a version must satisfy to
// be able to be selected.
type Criteria struct {
	// Below are fields for internal use only. Specifically used for
	// constraint parsing and checking.
	c          *semver.Constraints
	prerelease string

	once sync.Once

	// Constraint is a semantic versioning constraint that the version
	// must satisfy.
	//
	// Example: ">=1.0.0 <2.0.0"
	Constraint string

	// Branch is the branch that the version must point to. This
	// constraint will only be satisfied if the branch currently points to
	// the commit being considered.
	//
	// If a branch is provided, it will always be used over other
	// versions. For this reason, top-level modules should only ever use
	// branches.
	Branch string
}

// Parse parses the criteria's constraint into a semver constraint. If
// the constraint is already parsed, this is a no-op.
func (c *Criteria) Parse() error {
	var err error
	c.once.Do(func() {
		if c.Constraint == "" {
			// No constraint, no need to parse.
			return
		}

		if strings.Contains(c.Constraint, "||") || strings.Contains(c.Constraint, "&&") {
			// We don't support complex constraints.
			err = fmt.Errorf("complex constraints are not supported")
			return
		}

		// Create a "version" from the constraint
		cv := constRexp.ReplaceAllString(c.Constraint, "")

		// Attempt to parse the constraint as a version for detecting
		// per-release versions.
		vc, err := semver.NewVersion(cv)
		if err == nil {
			c.prerelease = strings.Split(vc.Prerelease(), ".")[0]
		}

		c.c, err = semver.NewConstraint(c.Constraint)
		if err != nil {
			return
		}
	})

	return err
}

// Check returns true if the version satisfies the criteria. If a
// prerelease is included then the provided criteria will be mutated to
// support pre-releases as well as ensure that the prerelease string
// matches the provided version. If a branch is provided, then the
// criteria will always be satisfied unless the criteria is looking for
// a specific branch, in which case it will be satisfied only if the
// branches match.
func (c *Criteria) Check(v *Version, prerelease, branch string) bool {
	if c.Branch != "" && v.Branch == c.Branch {
		return true
	}

	// Looking for a specific branch, but we're not asking for a branch,
	// so return success because we cannot compare these versions.
	if branch != "" && c.Branch == "" {
		return true
	}

	if c.c != nil && v.sv != nil {
		if c.prerelease != "" && c.prerelease != prerelease {
			// The provided criteria has a pre-release version, but the
			// version we're checking against does not match. This means
			// that we should not consider this version.
			return false
		}

		// If we're eligible for pre-releases but our constraint doesn't
		// allow for them, then we need to change our constraint to allow
		// for pre-releases.
		if prerelease != "" && c.prerelease == "" {
			// We need to add the pre-release to the constraint.
			c.Constraint = fmt.Sprintf("%s-%s", c.Constraint, prerelease)

			// TODO(jaredallard): Better error handling and location for this logic since
			// doing this on every call is pretty awful and inefficient.
			var err error
			c.c, err = semver.NewConstraint(c.Constraint)
			if err != nil {
				// This should never happen since we've already parsed
				// the constraint once.
				panic(fmt.Sprintf("failed to parse constraint: %v", err))
			}
			c.prerelease = prerelease
		}

		return c.c.Check(v.sv)
	}

	// Otherwise, doesn't match.
	return false
}

// Equal returns true if the criteria is equal to the other criteria.
func (c *Criteria) Equal(other *Criteria) bool {
	// If either is nil, they must both be nil.
	if c == nil || other == nil {
		return c == other
	}

	// Otherwise, check all fields.
	return c.Constraint == other.Constraint && c.Branch == other.Branch
}

// String returns a user-friendly representation of the criteria.
func (c *Criteria) String() string {
	if c.Branch != "" {
		return fmt.Sprintf("branch %s", c.Branch)
	}

	return c.Constraint
}
