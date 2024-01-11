//go:build !goverter

package db

func init() {
	c = &GameConverter{}
}
