package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
	"strconv"
)

func (p *Parser) parserNotFound(t token.TokenType) {
	p.addError(
		fmt.Sprintf("no prefix parser for %s has been found", t))
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	val, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer: %s",
			p.curToken.Literal, err.Error())
		p.addError(msg)
		return nil
	}
	return &ast.IntegerLiteral{
		Token: p.curToken,
		Value: val,
	}
}
