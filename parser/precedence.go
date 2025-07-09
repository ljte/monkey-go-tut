package parser

import "monkey/token"

type OperatorPrecedence int

const (
	_ OperatorPrecedence = iota

	PRECEDENCE_LOWEST
	PRECEDENCE_EQUALS
	PRECEDENCE_LESSGREATER
	PRECEDENCE_SUM
	PRECEDENCE_PRODUCT
	PRECEDENCE_PREFIX
	PRECEDENCE_CALL
)

var Precedences = map[token.TokenType]OperatorPrecedence{
	token.EQUALS:     PRECEDENCE_EQUALS,
	token.NOT_EQUALS: PRECEDENCE_EQUALS,
	token.LT:         PRECEDENCE_LESSGREATER,
	token.GT:         PRECEDENCE_LESSGREATER,
	token.PLUS:       PRECEDENCE_SUM,
	token.MINUS:      PRECEDENCE_SUM,
	token.SLASH:      PRECEDENCE_PRODUCT,
	token.ASTERISK:   PRECEDENCE_PRODUCT,
}

func DerivePrecedence(tt token.TokenType) OperatorPrecedence {
	if p, ok := Precedences[tt]; ok {
		return p
	}
	return PRECEDENCE_LOWEST
}
