package main

import (
	"fmt"
	"os"

	"github.com/agustin-del-pino/nadia-lang/nadia/binder"
	"github.com/agustin-del-pino/nadia-lang/nadia/parser"
	"github.com/agustin-del-pino/nadia-lang/nadia/transpiler"
)

func main() {
	code, cErr := os.ReadFile("./main.nad")
	if cErr != nil {
		fmt.Fprintln(os.Stderr, cErr)
		return
	}
	src, sErr := parser.Parse(code)
	if sErr != nil {
		fmt.Fprintln(os.Stderr, sErr)
		return
	}
	bErr := binder.Bind("./", "../nad/lib", src, nil)
	if bErr != nil {
		fmt.Fprintln(os.Stderr, bErr)
		return
	}
	bdl := transpiler.Transpile(src)
	out := transpiler.Format(bdl, 2)
	fmt.Println(string(out))
}
