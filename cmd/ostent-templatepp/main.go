package main

import (
	"github.com/ostrost/ostent/templateutil/templatefunc"
	"github.com/ostrost/ostent/templateutil/templatepipe"
)

func main() {
	templatepipe.Main(
		templatefunc.FuncMapHTML(),
		templatefunc.FuncMapJSXL(),
	)
}
