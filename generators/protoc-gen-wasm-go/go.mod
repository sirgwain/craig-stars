module github.com/sirgwain/craig-stars/generators/protoc-gen-wasm-go

go 1.27.1

replace github.com/sirgwain/craig-stars/proto-wasm => ../../proto-wasm

require (
	github.com/sirgwain/craig-stars/proto-wasm v0.0.0-20261006155015-1544c51b8442
	google.golang.org/protobuf v1.36.12
)
