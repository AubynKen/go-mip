# MIP in Go demo project

This project implements a simple mixed-integer programming (MIP) model in Go using the OR-Tools library.
We use the C++ api of OR-Tools through cgo.

## Structure
`bridge`: the C++ bridging layer code that interfaces with the OR-Tools library.
`mip`: the Go wrapper code that interfaces with the C++ bridging layer.
`examples`: example MIP models that can be solved using the Go wrapper code. (For now the SD-wan example is in main.go)

## Build
