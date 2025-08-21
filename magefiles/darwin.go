//go:build darwin && !wasi && !wasm

package main

const binary_name = "craig-stars"
const ldflags = "-ldflags=-s -w"
