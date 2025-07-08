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
let x = 10;
let t = 13123;
let ad = 1231;
`
	p := parser.New(lexer.New(input))

	program := p.ParseProgram()
	fmt.Println(program.String())
	for _, stmt := range program.Statements {
		fmt.Println(stmt.TokenLiteral())
	}

	m := map[string]int{}

	fmt.Println(m)

	// l := lexer.New("!=")
	// next := l.NextToken()
	// fmt.Printf("%v %q\n", next.Type, next.Literal)
	// fmt.Println(l.Next
	// fmt.Println(l.NextToken())
}
