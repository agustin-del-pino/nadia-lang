package main

import (
	"os"
	"path/filepath"

	"github.com/agustin-del-pino/nadia-lang/nad"
)

func cmdInit(c cli) error {
	err := os.MkdirAll(c.bin, 0777)
	if err != nil {
		return err
	}

	return nad.ReadFiles(func(n string, b []byte) error {
		return os.WriteFile(filepath.Join(c.bin, n), b, 0777)
	})
}
