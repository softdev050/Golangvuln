// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package osv implements the OSV shared vulnerability
// format, as defined by https://ossf.github.io/osv-schema.
//
// As this package is intended for use with the Go vulnerability
// database, only the subset of features which are used by that
// database are implemented (for instance, only the SEMVER affected
// range type is implemented).
package osv

import (
	"time"

	"golang.org/x/mod/semver"
	isem "golang.org/x/vuln/internal/semver"
)

type AffectsRangeType string

const (
	TypeUnspecified AffectsRangeType = "UNSPECIFIED"
	TypeGit         AffectsRangeType = "GIT"
	TypeSemver      AffectsRangeType = "SEMVER"
)

type Ecosystem string

const GoEcosystem Ecosystem = "Go"

type Package struct {
	Name      string    `json:"name"`
	Ecosystem Ecosystem `json:"ecosystem"`
}

type RangeEvent struct {
	Introduced string `json:"introduced,omitempty"`
	Fixed      string `json:"fixed,omitempty"`
}

type AffectsRange struct {
	Type   AffectsRangeType `json:"type"`
	Events []RangeEvent     `json:"events"`
}

func (ar AffectsRange) containsSemver(v string) bool {
	if ar.Type != TypeSemver {
		return false
	}
	if len(ar.Events) == 0 {
		return true
	}

	// Strip and then add the semver prefix so we can support bare versions,
	// versions prefixed with 'v', and versions prefixed with 'go'.
	v = isem.CanonicalizeSemverPrefix(v)

	var affected bool
	for _, e := range ar.Events {
		if !affected && e.Introduced != "" {
			affected = e.Introduced == "0" || semver.Compare(v, isem.CanonicalizeSemverPrefix(e.Introduced)) >= 0
		} else if e.Fixed != "" {
			affected = semver.Compare(v, isem.CanonicalizeSemverPrefix(e.Fixed)) < 0
		}
	}

	return affected
}

type Affects []AffectsRange

func (a Affects) AffectsSemver(v string) bool {
	if len(a) == 0 {
		// No ranges implies all versions are affected
		return true
	}
	var semverRangePresent bool
	for _, r := range a {
		if r.Type != TypeSemver {
			continue
		}
		semverRangePresent = true
		if r.containsSemver(v) {
			return true
		}
	}
	// If there were no semver ranges present we
	// assume that all semvers are affected, similarly
	// to how to we assume all semvers are affected
	// if there are no ranges at all.
	return !semverRangePresent
}

type Reference struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type Affected struct {
	Package           Package           `json:"package"`
	Ranges            Affects           `json:"ranges,omitempty"`
	DatabaseSpecific  DatabaseSpecific  `json:"database_specific"`
	EcosystemSpecific EcosystemSpecific `json:"ecosystem_specific"`
}

type DatabaseSpecific struct {
	URL string `json:"url"`
}

// EcosytemSpecificImport contains additional information about an affected package.
type EcosystemSpecificImport struct {
	// Path is the package import path.
	Path string `json:"path,omitempty"`

	// Symbols is the collection of functions and methods names affected by
	// this vulnerability. Methods are listed as <recv>.<method>.
	//
	// If included, only programs which use these symbols will be marked as
	// vulnerable. If omitted, any program which imports this module will be
	// marked vulnerable.
	//
	// These should be the symbols initially detected or identified in the CVE
	// or other source.
	Symbols []string `json:"symbols,omitempty"`
}

// EcosystemSpecific contains additional information about the vulnerability
// for the Go ecosystem.
type EcosystemSpecific struct {
	// Imports is the list of affected packages within the module.
	Imports []EcosystemSpecificImport `json:"imports,omitempty"`

	// GOOS is the execution operating system where the symbols appear, if
	// known.
	//
	// At the moment, this information is not provided by the Go
	// vulnerability database.
	GOOS []string `json:"goos,omitempty"`

	// GOARCH specifies the execution architecture where the symbols appear, if
	// known.
	//
	// At the moment, this information is not provided by the Go
	// vulnerability database.
	GOARCH []string `json:"goarch,omitempty"`
}

// Entry represents a OSV style JSON vulnerability database
// entry
type Entry struct {
	ID         string      `json:"id"`
	Published  time.Time   `json:"published"`
	Modified   time.Time   `json:"modified"`
	Withdrawn  *time.Time  `json:"withdrawn,omitempty"`
	Aliases    []string    `json:"aliases,omitempty"`
	Details    string      `json:"details"`
	Affected   []Affected  `json:"affected"`
	References []Reference `json:"references,omitempty"`
}
