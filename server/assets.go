package server

import "embed"

// assets from frontend/build
// this is Set in main.go because we can't
// embed ../ relative paths.
var assets embed.FS

func SetAssets(value embed.FS) {
	assets = value
}

func GetAssets() embed.FS {
	return assets
}
