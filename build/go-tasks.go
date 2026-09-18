package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/goyek/goyek/v3"
	"github.com/goyek/x/cmd"
)

var goGenerate = goyek.Define(goyek.Task{
	Name:  "go-generate",
	Usage: "go generate",
	Action: func(a *goyek.A) {
		cmd.Exec(a, "go generate ./...")
	},
})

var goTest = goyek.Define(goyek.Task{
	Name:  "go-test",
	Usage: "go test",
	Action: func(a *goyek.A) {
		out := filepath.Join("bin", "coverage.out")
		html := filepath.Join("bin", "coverage.html")
		var pkgs strings.Builder
		if !cmd.Exec(a, "go list -f '{{ .ImportPath }}' ./...", cmd.Stdout(&pkgs)) {
			return
		}

		// Exclude this package from the code coverage package list
		info, ok := debug.ReadBuildInfo()
		if !ok {
			a.Fatal("could not read `./build` build info")
			return
		}
		packageList := strings.Split(strings.TrimSpace(pkgs.String()), "\n")
		packageList = slices.DeleteFunc(packageList, func(s string) bool {
			return strings.HasSuffix(s, info.Path)
		})

		packages := strings.Join(append(packageList, "."), ",")
		if !cmd.Exec(a,
			fmt.Sprintf("go test -race -covermode=atomic -coverprofile=%q -coverpkg=%q ./...", out, packages),
		) {
			return
		}
		cmd.Exec(a, fmt.Sprintf("go tool cover -html=%q -o %q", out, html))
	},
	Deps: []*goyek.DefinedTask{mkdirBin},
})

var goBuild = goyek.Define(goyek.Task{
	Name:  "go-build",
	Usage: "go build",
	Action: func(a *goyek.A) {
		var stderr strings.Builder
		if !cmd.Exec(a, `go build -trimpath -ldflags="-dumpdep -s -w" -o bin/ .`, cmd.Stderr(&stderr)) {
			return
		}

		err := os.WriteFile(filepath.Join("bin", "deps.txt"), []byte(stderr.String()), 0600)
		if err != nil {
			a.Fatal(err)
		}
	},
	Deps: []*goyek.DefinedTask{mkdirBin},
})

var goModTidyDiff = goyek.Define(goyek.Task{
	Name:  "go-mod-tidy",
	Usage: "go mod tidy",
	Action: func(a *goyek.A) {
		cmd.Exec(a, "go mod tidy -diff")
	},
})
