module github.com/sirgwain/craig-stars/cs

go 1.24.2

replace github.com/sirgwain/craig-stars/test => ../test

require (
	github.com/rs/zerolog v1.34.0
	github.com/sirgwain/craig-stars/test v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.10.0
	golang.org/x/crypto v0.41.0
	golang.org/x/exp v0.0.0-20250819193227-8b4c13bb791b
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/nsf/jsondiff v0.0.0-20230430225905-43f6cf3098c1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	golang.org/x/sys v0.35.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
