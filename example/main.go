package main

import (
	"fmt"
	"os"

	"github.com/agustin-del-pino/nadia-lang/nadia/parser"
	"github.com/agustin-del-pino/nadia-lang/nadia/transpiler"
)

const code = `
type int -> 0 -> "int" {
    (1, -), (2, +), (2, -), 
    (2, *), (2, ^), (2, mod), 
    (2, ==), (2, <), (2, >), 
    (2, !=), (2, >=), (2, <=)
}

type float -> 0.0 -> "float" {
    (1, -), (2, +), (2, -), 
    (2, *), (2, ^), (2, mod), 
    (2, ==), (2, <), (2, >), 
    (2, !=), (2, >=), (2, <=)
}

type bool -> False -> "bool" {
    (1, not), (2, ==), (2, !=),
    (2, and), (2, or)
}


type signal -> Low -> "int" {
    (1, not), (2, ==), (2, !=)
}

type string -> "" -> "String" {
    (2, +), (2, ==), (2, !=)
}

def Low signal  = "LOW"
def High signal = "HIGH"
def True bool   = "true"
def False bool  = "false"

def ref(any) any = "(&$1)"
def ptr(any) any = "(*$1)"

def event_of(any) any = "event_$1"
def trigger_of(any) bool = "ev_$1_trigger()"


func main() {
	for i range 0, 10 {
		if i mod 2 == 0 {
			pass
		}
		
		if i == 7 {
			stop
		}
	}
}
`

func main() {
	src, sErr := parser.Parse([]byte(code))
	if sErr != nil {
		fmt.Fprint(os.Stderr, sErr)
	}

	bdl := transpiler.Transpile(src)
	out := transpiler.Format(bdl, 2)
	fmt.Println(string(out))
}
