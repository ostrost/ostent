//+build none

package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, flag.Arg(0), nil, 0)
	if err != nil {
		panic(err)
	}
	ast.Print(fset, f)
}
