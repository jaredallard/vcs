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

// Package fileinfo provides a simple implementation of [os.FileInfo]
// for remote files.
package fileinfo

import (
	"io/fs"
	"os"
	"time"
)

// _ ensures that [fileInfo] implements [os.FileInfo].
var _ os.FileInfo = &File{}

// File implements [os.FileInfo] for a given release asset. Given
// that these are remote files, not all fields are implemented:
//
//   - ModTime: If supported by the underlying VCS provider, created at
//     will be used instead.
//   - IsDir: Always returns "false".
//   - Mode: Always returns 0o600.
//   - Sys: Returns the underlying struct used to create this, if set by
//     the VCS provider. This CAN return "nil".
type File struct {
	sys     any
	modTime time.Time
	name    string
	size    int64
}

// New creates a new [File] instance with the given parameters.
func New(name string, size int64, modTime time.Time, sys any) *File {
	return &File{
		sys:     sys,
		modTime: modTime,
		name:    name,
		size:    size,
	}
}

// IsDir implements [os.FileInfo], see [fileInfo] and the previously
// mentioned interface for more information.
func (f *File) IsDir() bool {
	return false
}

// ModTime implements [os.FileInfo], see [fileInfo] and the previously
// mentioned interface for more information.
func (f *File) ModTime() time.Time {
	return f.modTime
}

// Mode implements [os.FileInefo], see [fileInfo] and the previously
// mentioned interface for more information.
func (f *File) Mode() fs.FileMode {
	return fs.FileMode(0o600)
}

// Name implements [os.FileInfo], see [fileInfo] and the previously
// mentioned interface for more information.
func (f *File) Name() string {
	return f.name
}

// Size implements [os.FileInfo], see [fileInfo] and the previously
// mentioned interface for more information.
func (f *File) Size() int64 {
	return f.size
}

// Sys implements [os.FileInfo], see [fileInfo] and the previously
// mentioned interface for more information.
func (f *File) Sys() any {
	return f.sys
}
