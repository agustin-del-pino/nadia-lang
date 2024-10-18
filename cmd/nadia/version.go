package main

import "fmt"

const version = "v0.1.0"

func cmdVersion(c cli) error {
	fmt.Println("Nadia-Lang ", version)
	return nil
}
