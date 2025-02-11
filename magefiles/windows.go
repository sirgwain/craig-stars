//go:build windows

package main

const binary_name = "craig-stars.exe"
const ldflags = "-ldflags=-s -w -extldflags '-static'"
