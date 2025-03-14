package lexer

import (
	"testing"

	"monkey/token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LCURLY, "{"},
		{token.RCURLY, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		token := l.NextToken()

		if token.Type != tt.expectedType {
			t.Fatalf("tests [%d] failed, expected tokentype: %s, but got: %s",
				i, tt.expectedType, token.Type)
		}

		if token.Literal != tt.expectedLiteral {
			t.Fatalf("tests [%d] failed, expected literal: %s, but got: %s",
				i, tt.expectedLiteral, token.Literal)
		}
	}
}
