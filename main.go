package main

import (
	"fmt"
	"monkey/lexer"
	"monkey/parser"
)

func main() {
	// user, err := user.Current()
	// if err != nil {
	// 	panic(err)
	// }
	// repl := repl.New(os.Stdin, os.Stdout)

	// if err := repl.Loop(); err != nil {
	// 	panic(err)
	// }
	//
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	p := parser.New(lexer.New(input))

	program := p.ParseProgram()
	for _, stmt := range program.Statements {
		fmt.Println(stmt.TokenLiteral())
	}

	// l := lexer.New("!=")
	// next := l.NextToken()
	// fmt.Printf("%v %q\n", next.Type, next.Literal)
	// fmt.Println(l.Next
	// fmt.Println(l.NextToken())
}
