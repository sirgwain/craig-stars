module github.com/sirgwain/craig-stars/generators/protoc-gen-wasm-go

go 1.24.2

replace github.com/sirgwain/craig-stars/proto-wasm => ../../proto-wasm

require (
	github.com/sirgwain/craig-stars/proto-wasm v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.36.7
)
