// Package version exposes build-time metadata for the craig-stars binary.
package version

import "time"

// Semver is the semantic-release version injected at build time.
var Semver = "0.0.0-develop"

// Commit is the source commit injected at build time.
var Commit = "local"

// BuildTime is the release build time injected at build time.
var BuildTime = time.Now().Format(time.RFC3339)
