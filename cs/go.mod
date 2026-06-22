module github.com/sirgwain/craig-stars/cs

go 1.26.1

replace github.com/sirgwain/craig-stars/test => ../test

require (
	github.com/sirgwain/craig-stars/test v0.0.0-20260621133813-709c3b17b2dc
	github.com/stretchr/testify v1.11.1
	golang.org/x/exp v0.0.0-20260611194520-c48552f49976
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/nsf/jsondiff v0.0.0-20260207060731-8e8d90c4c0ac // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
