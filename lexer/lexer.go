package lexer

import "monkey/token"

type Lexer struct {
	input       string
	pos         int
	readPos     int
	currentChar byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input, 0, 0, 0}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.currentChar {
	case '=':
		tok = newtoken(token.ASSIGN, l.currentChar)

	case '+':
		tok = newtoken(token.PLUS, l.currentChar)

	case ',':
		tok = newtoken(token.COMMA, l.currentChar)

	case ';':
		tok = newtoken(token.SEMICOLON, l.currentChar)

	case '(':
		tok = newtoken(token.LPAREN, l.currentChar)

	case ')':
		tok = newtoken(token.RPAREN, l.currentChar)

	case '{':
		tok = newtoken(token.LCURLY, l.currentChar)

	case '}':
		tok = newtoken(token.RCURLY, l.currentChar)

	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}
	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.currentChar = 0
	} else {
		l.currentChar = l.input[l.readPos]
	}

	l.pos = l.readPos
	l.readPos += 1
}

func newtoken(ttype token.TokenType, ch byte) token.Token {
	return token.Token{ttype, string(ch)}
}
