module github.com/sirgwain/craig-stars/wasm

go 1.24.2

replace github.com/sirgwain/craig-stars => ../

replace github.com/sirgwain/craig-stars/cs => ../cs

replace github.com/sirgwain/craig-stars/test => ../test

replace github.com/sirgwain/craig-stars/proto-wasm => ../proto-wasm

require (
	github.com/sirgwain/craig-stars v0.0.0-00010101000000-000000000000
	github.com/sirgwain/craig-stars/cs v0.0.0-00010101000000-000000000000
)

require (
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/rs/zerolog v1.34.0
	google.golang.org/protobuf v1.36.8 // indirect
)

require (
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirgwain/craig-stars/proto-wasm v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/exp v0.0.0-20250819193227-8b4c13bb791b // indirect
	golang.org/x/sys v0.35.0 // indirect
)
