module github.com/sirgwain/craig-stars/generators/protoc-gen-wasm-go

go 1.24.2

replace github.com/sirgwain/craig-stars/proto-wasm => ../../proto-wasm

require (
	github.com/sirgwain/craig-stars/proto-wasm v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.36.7
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
)
