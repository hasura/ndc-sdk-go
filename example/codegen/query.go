package main

import (
	"context"
)

// A hello argument
type HelloArguments struct {
	Num int
	Str string
}

// A hello result
type HelloResult struct {
	Num int
	Str string
}

// Send hello message
//
// @procedure hello
func Hello(ctx context.Context, state *State, arguments *HelloArguments) (*HelloResult, error) {
	return &HelloResult{
		Num: 1,
		Str: "world",
	}, nil
}
