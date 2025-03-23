package main

var Aliases = map[string]any{
	"dev":          Launch,
	"dev_frontend": Launch_Frontend,
	"dev_backend":  Launch_Backend,
	"copy_wasm":    Copy_Wasm_Exec,
	"test_backend": Test_Golang,
}
