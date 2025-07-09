package main

import (
	"fmt"
	"monkey/ast"
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
5;
`
	p := parser.New(lexer.New(input))

	program := p.ParseProgram()
	for _, stmt := range program.Statements {
		integ, _ := stmt.(*ast.ExpressionStatement)
		fmt.Println(integ.Expression.String())
	}

	// l := lexer.New("!=")
	// next := l.NextToken()
	// fmt.Printf("%v %q\n", next.Type, next.Literal)
	// fmt.Println(l.Next
	// fmt.Println(l.NextToken())
}
