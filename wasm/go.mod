module github.com/sirgwain/craig-stars/wasm

go 1.25

replace github.com/sirgwain/craig-stars => ../

replace github.com/sirgwain/craig-stars/cs => ../cs

replace github.com/sirgwain/craig-stars/test => ../test

replace github.com/sirgwain/craig-stars/proto-wasm => ../proto-wasm

require (
	github.com/sirgwain/craig-stars v0.0.0-00010101000000-000000000000
	github.com/sirgwain/craig-stars/cs v0.0.0-00010101000000-000000000000
)

require (
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	google.golang.org/protobuf v1.36.8 // indirect
)

require (
	github.com/sirgwain/craig-stars/proto-wasm v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/exp v0.0.0-20250819193227-8b4c13bb791b // indirect
)
