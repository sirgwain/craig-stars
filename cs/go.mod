module github.com/sirgwain/craig-stars/cs

go 1.25

replace github.com/sirgwain/craig-stars/test => ../test

require (
	github.com/sirgwain/craig-stars/test v0.0.0-20250827144401-54aa2fd68ef6
	github.com/stretchr/testify v1.11.1
	golang.org/x/exp v0.0.0-20250819193227-8b4c13bb791b
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/nsf/jsondiff v0.0.0-20230430225905-43f6cf3098c1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
