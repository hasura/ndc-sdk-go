package main

import (
	"context"
)

// A hello argument
type HelloArguments struct {
	Num int    `json:"num"`
	Str string `json:"str"`
}

// A hello result
type HelloResult struct {
	Num int    `json:"num"`
	Str string `json:"str"`
}

// Send hello message
//
// @function hello
func Hello(ctx context.Context, state *State, arguments *HelloArguments) (*HelloResult, error) {
	return &HelloResult{
		Num: 1,
		Str: "world",
	}, nil
}

// A create author argument
type CreateAuthorArguments struct {
	Name string `json:"name"`
}

// A create author result
type CreateAuthorResult struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// A procedure example
//
// @procedure
func CreateAuthor(ctx context.Context, state *State, arguments *CreateAuthorArguments) (*CreateAuthorResult, error) {
	return &CreateAuthorResult{
		ID:   1,
		Name: arguments.Name,
	}, nil
}
