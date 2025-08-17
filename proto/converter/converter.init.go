//go:build !goverter

package converter

func init() {
	C = &ProtoConverter{}
}
