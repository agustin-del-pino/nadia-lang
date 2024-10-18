package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var commands = map[string]func(cli) error{
	"init":    cmdInit,
	"build":   cmdBuild,
	"help":    cmdHelp,
	"version": cmdVersion,
}

var flags = map[string]struct{}{
	"-o":           {},
	"--out":        {},
	"-n":           {},
	"--nadia-path": {},
}

type cli struct {
	flags   map[string]string
	command string
	param   string
	cwd     string
	bin     string
}

func (c cli) GetFlag(fls ...string) (string, bool) {
	for _, f := range fls {
		if v, ok := c.flags[f]; ok {
			return v, ok
		}
	}

	return "", false
}

func getUserInput(c *cli) error {
	args := os.Args[1:]
	c.command = args[0]
	if _, ok := commands[c.command]; !ok {
		return fmt.Errorf("invalid command: %s", c.command)
	}

	sArgs := args[1:]
	if len(sArgs) == 0 {
		return nil
	}

	c.param = sArgs[0]

	flg := sArgs[1:]

	if len(flg) == 0 {
		return nil
	}

	c.flags = make(map[string]string)
	var f string

	for i := 0; i < len(flg); {
		if f == "" {
			if _, ok := flags[flg[i]]; !ok {
				i++
				continue
			}
			f = flg[i]
			c.flags[f] = ""
			i++
		} else {
			if flg[i][0] == '-' {
				f = ""
				continue
			}
			c.flags[f] = flg[i]
			i++
		}
	}
	return nil
}

func main() {
	cwd, err := os.Getwd()

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	bin, bErr := os.Executable()
	if bErr != nil {
		fmt.Fprintln(os.Stderr, bErr.Error())
	}

	var c cli
	c.cwd = cwd
	c.bin = filepath.Join(filepath.Dir(bin), "nad", "lib")

	err = getUserInput(&c)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	err = commands[c.command](c)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
