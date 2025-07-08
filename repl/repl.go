package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
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

	for {
		fmt.Fprint(r.out, r.prompt)
		if !scanner.Scan() {
			return nil
		}
		l := lexer.New(scanner.Text())

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(r.out, "%+v\n", tok)
		}
	}

}
