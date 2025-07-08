package main

import (
	"monkey/repl"
	"os"
)

func main() {
	// user, err := user.Current()
	// if err != nil {
	// 	panic(err)
	// }
	repl := repl.New(os.Stdin, os.Stdout)

	if err := repl.Loop(); err != nil {
		panic(err)
	}

	// l := lexer.New("!=")
	// next := l.NextToken()
	// fmt.Printf("%v %q\n", next.Type, next.Literal)
	// fmt.Println(l.Next
	// fmt.Println(l.NextToken())
}
