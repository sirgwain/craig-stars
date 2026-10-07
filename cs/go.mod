module github.com/sirgwain/craig-stars/cs

go 1.27.1

replace github.com/sirgwain/craig-stars/test => ../test

require (
	github.com/sirgwain/craig-stars/test v0.0.0-20261006155015-1544c51b8442
	github.com/stretchr/testify v1.12.1
	golang.org/x/exp v0.0.0-20261005173118-76772065c9b0
)

require (
	github.com/nsf/jsondiff v0.0.0-20260207060731-8e8d90c4c0ac // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
)
