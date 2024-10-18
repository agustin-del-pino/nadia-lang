package main

import (
	"os"
	"strings"

	"github.com/agustin-del-pino/nadia-lang/nadia/binder"
	"github.com/agustin-del-pino/nadia-lang/nadia/parser"
	"github.com/agustin-del-pino/nadia-lang/nadia/transpiler"
)

func cmdBuild(c cli) error {
	b, bErr := os.ReadFile(c.param)
	if bErr != nil {
		return bErr
	}

	src, sErr := parser.Parse(b)
	if sErr != nil {
		return sErr
	}
	var ext string
	if v, ok := c.GetFlag("-n", "--nadia-path"); ok {
		ext = v
	} else {
		ext = c.bin
	}
	err := binder.Bind(c.cwd, ext, src, nil)
	if err != nil {
		return err
	}
	bdl := transpiler.Transpile(src)
	out := transpiler.Format(bdl, 2)

	var pth string

	if v, ok := c.GetFlag("-o", "--out"); ok {
		pth = v
	} else {
		pth = strings.ReplaceAll(c.param, ".nad", ".ino")
	}

	return os.WriteFile(pth, out, 0777)
}
