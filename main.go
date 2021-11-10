package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Pkg struct {
	Deps        map[string]bool
	ReverseDeps map[string]bool
}

type TrailEntry struct {
	PkgNames []string
	PkgNext  int
}

func main() {
	flag.Parse()

	pkgs := make(map[string]Pkg)

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		entries := strings.Split(s.Text(), " ")
		node := entries[0]
		dep := entries[1]

		nodePkg := getPkg(pkgs, node)
		depPkg := getPkg(pkgs, dep)

		nodePkg.Deps[dep] = true
		depPkg.ReverseDeps[node] = true
	}

	printed := make(map[string]bool)
	trail := []*TrailEntry{{
		PkgNames: []string{flag.Arg(0)},
		PkgNext:  0,
	}}

	for len(trail) > 0 {
		current := trail[len(trail)-1]
		if current.PkgNext >= len(current.PkgNames) {
			trail = trail[:len(trail)-1]
			continue
		}

		nextName := current.PkgNames[current.PkgNext]
		current.PkgNext++

		suffix := ""
		if printed[nextName] {
			suffix = " ^^"
		}
		fmt.Printf("%s%s%s\n", strings.Repeat("|   ", len(trail)-1), nextName, suffix)

		if printed[nextName] {
			continue
		}

		nextPkg := pkgs[nextName]
		nextPkgNames := make([]string, len(nextPkg.ReverseDeps))

		i := 0
		for k := range nextPkg.ReverseDeps {
			nextPkgNames[i] = k
			i++
		}
		sort.Strings(nextPkgNames)

		trail = append(trail, &TrailEntry{
			PkgNames: nextPkgNames,
			PkgNext:  0,
		})

		printed[nextName] = true
	}
}

func getPkg(pkgs map[string]Pkg, name string) Pkg {
	pkg, ok := pkgs[name]
	if !ok {
		pkg = Pkg{
			Deps:        make(map[string]bool),
			ReverseDeps: make(map[string]bool),
		}
		pkgs[name] = pkg
	}

	return pkg
}
