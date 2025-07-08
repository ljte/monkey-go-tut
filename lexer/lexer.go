package lexer

import (
	"monkey/token"
)

type Lexer struct {
	input       string
	pos         int
	readPos     int
	currentChar byte
}

func New(input string) *Lexer {
	l := &Lexer{input, 0, 0, 0}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespaces()

	switch l.currentChar {
	case '=':
		nextChar := l.peek()
		if nextChar == '=' {
			tok = newtoken(token.EQUALS, string(l.currentChar)+string(nextChar))
			l.readChar()
		} else {
			tok = newtoken(token.ASSIGN, string(l.currentChar))
		}

	case '+':
		tok = newtoken(token.PLUS, string(l.currentChar))

	case ',':
		tok = newtoken(token.COMMA, string(l.currentChar))

	case ';':
		tok = newtoken(token.SEMICOLON, string(l.currentChar))

	case '(':
		tok = newtoken(token.LPAREN, string(l.currentChar))

	case ')':
		tok = newtoken(token.RPAREN, string(l.currentChar))

	case '{':
		tok = newtoken(token.LCURLY, string(l.currentChar))

	case '}':
		tok = newtoken(token.RCURLY, string(l.currentChar))

	case '!':
		nextChar := l.peek()
		if nextChar == '=' {
			tok = newtoken(token.NOT_EQUALS, string(l.currentChar)+string(nextChar))
			l.readChar()
		} else {
			tok = newtoken(token.BANG, string(l.currentChar))
		}

	case '*':
		tok = newtoken(token.ASTERISK, string(l.currentChar))

	case '/':
		tok = newtoken(token.SLASH, string(l.currentChar))

	case '-':
		tok = newtoken(token.MINUS, string(l.currentChar))

	case '<':
		tok = newtoken(token.LT, string(l.currentChar))

	case '>':
		tok = newtoken(token.GT, string(l.currentChar))

	case 0:
		tok = newtoken(token.EOF, "")

	default:
		// if isWhitespace(l.currentChar) {
		// 	l.skipWhitespaces()
		// 	return l.NextToken()
		// }

		if isLetter(l.currentChar) {
			ident := l.readIdent()
			return newtoken(identOrKeyword(ident), ident)
		}

		if isInt(l.currentChar) {
			num := l.readInt()
			// if err != nil {
			// 	panic(err)
			// }
			return newtoken(token.INT, num)
		}
		tok = newtoken(token.ILLEGAL, string(l.currentChar))
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

func (l *Lexer) skipWhitespaces() {
	for isWhitespace(l.currentChar) {
		l.readChar()
	}
}

func (l *Lexer) readIdent() string {
	pos := l.pos
	for isLetter(l.currentChar) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readInt() string {
	pos := l.pos
	for isInt(l.currentChar) {
		l.readChar()
	}
	return l.input[pos:l.pos]
	// n, err := strconv.Atoi(num)
	// if err != nil {
	// 	return -1, fmt.Errorf("expected integer, got %v: %s", num, err)
	// }
	// return n, nil
}

func (l *Lexer) peek() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func identOrKeyword(ident string) token.TokenType {
	switch ident {
	case "let":
		return token.LET
	case "fn":
		return token.FUNCTION
	case "if":
		return token.IF
	case "else":
		return token.ELSE
	case "return":
		return token.RETURN
	case "true":
		return token.TRUE
	case "false":
		return token.FALSE
	default:
		return token.IDENTIFIER
	}
}

func newtoken(ttype token.TokenType, literal string) token.Token {
	return token.Token{Type: ttype, Literal: literal}
}

func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func isInt(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\n' || ch == '\t'
}
