//go:build !wasi && !wasm && !windows && !darwin

package main

const binary_name = "craig-stars"
const ldflags = "-ldflags=-s -w -extldflags '-static'"
