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

	// 	input := `
	// add(5, 6, 7)
	// `
	// 	p := parser.New(lexer.New(input))

	// 	program := p.ParseProgram()
	// 	// for _, stmt := range program.Statements {
	// 	// 	integ, _ := stmt.(*ast.ExpressionStatement)
	// 	// 	fmt.Println(integ.Expression.String())
	// 	// }
	// 	fmt.Println(program.String())

	// // l := lexer.New("!=")
	// // next := l.NextToken()
	// // fmt.Printf("%v %q\n", next.Type, next.Literal)
	// // fmt.Println(l.Next
	// // fmt.Println(l.NextToken())
}
