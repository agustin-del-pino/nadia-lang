package main

import "fmt"

const help = `
nadia <command> [param] ...flags


Commands

init:     initializes the internal dependencies. 
build:    builds the given nadia source file
help:     displays CLI's help
version:  displays nadia-lang's version


Flags

-o, --out:         sets the output filepath
-n, --nadia-path:  sets a temporally external lib path`

func cmdHelp(c cli) error {
	cmdVersion(c)
	fmt.Println(help)
	return nil
}
