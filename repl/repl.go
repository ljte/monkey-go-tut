package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/eval"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

type Repl struct {
	prompt string
	in     io.Reader
	out    io.Writer
}

func New(in io.Reader, out io.Writer) *Repl {
	return &Repl{prompt: ">> ", in: in, out: out}
}

func (r *Repl) Loop() error {
	scanner := bufio.NewScanner(r.in)
	env := object.NewEnvironment()

	for {
		fmt.Fprint(r.out, r.prompt)
		if !scanner.Scan() {
			return nil
		}
		p := parser.New(lexer.New(scanner.Text()))
		program := p.ParseProgram()
		if errors := p.Errors(); len(errors) != 0 {
			r.printErrors(errors)
			continue
		}
		res := eval.Eval(program, env)
		if res != nil {
			fmt.Fprintln(r.out, res.Inspect())
		}
	}
}

func (r *Repl) printErrors(errors []string) {
	for _, err := range errors {
		fmt.Fprintf(r.out, "\t%s\n", err)
	}
}
