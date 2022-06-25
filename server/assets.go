package server

import "embed"

var assets embed.FS

func SetAssets(value embed.FS) {
	assets = value
}

func GetAssets() embed.FS {
	return assets
}
