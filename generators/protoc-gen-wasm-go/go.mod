module github.com/sirgwain/craig-stars/generators/protoc-gen-wasm-go

go 1.25

replace github.com/sirgwain/craig-stars/proto-wasm => ../../proto-wasm

require (
	github.com/sirgwain/craig-stars/proto-wasm v0.0.0-20260621133813-709c3b17b2dc
	google.golang.org/protobuf v1.36.11
)
