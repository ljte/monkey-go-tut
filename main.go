package main

import (
	"monkey/repl"
	"os"
)

func main() {
	repl := repl.New(os.Stdin, os.Stdout)

	if err := repl.Loop(); err != nil {
		panic(err)
	}
}
