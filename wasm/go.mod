module github.com/sirgwain/craig-stars/wasm

go 1.25

replace github.com/sirgwain/craig-stars => ../

replace github.com/sirgwain/craig-stars/cs => ../cs

replace github.com/sirgwain/craig-stars/test => ../test

replace github.com/sirgwain/craig-stars/proto-wasm => ../proto-wasm

require (
	github.com/sirgwain/craig-stars v0.0.0-00010101000000-000000000000
	github.com/sirgwain/craig-stars/cs v0.0.0-20260621133813-709c3b17b2dc
)

require (
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

require (
	github.com/sirgwain/craig-stars/proto-wasm v0.0.0-20260621133813-709c3b17b2dc // indirect
	golang.org/x/exp v0.0.0-20260611194520-c48552f49976 // indirect
)
