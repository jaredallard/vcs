// Copyright (C) 2026 vcs contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

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
