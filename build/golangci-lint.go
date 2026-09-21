package main

import (
	"github.com/goyek/goyek/v3"
	"github.com/goyek/x/cmd"
)

var golangciLint = goyek.Define(goyek.Task{
	Name:  "lint",
	Usage: "lint",
	Action: func(a *goyek.A) {
		cmd.Exec(
			a,
			"go run -modfile=./tools/golangci-lint/go.mod github.com/golangci/golangci-lint/v2/cmd/golangci-lint run",
		)
	},
})

var _ = goyek.Define(goyek.Task{
	Name:  "lint-fix",
	Usage: "lint-fix",
	Action: func(a *goyek.A) {
		cmd.Exec(
			a,
			"go run -modfile=./tools/golangci-lint/go.mod github.com/golangci/golangci-lint/v2/cmd/golangci-lint run --fix",
		)
	},
})
