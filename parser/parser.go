package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type (
	prefixParser func() ast.Expression
	infixParser  func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string

	prefixParsers map[token.TokenType]prefixParser
	infixParsers  map[token.TokenType]infixParser
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:             l,
		errors:        []string{},
		prefixParsers: map[token.TokenType]prefixParser{},
		infixParsers:  map[token.TokenType]infixParser{},
	}

	p.registerPrefixParser(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefixParser(token.INT, p.parseIntegerLiteral)
	p.registerPrefixParser(token.BANG, p.parsePrefixExpression)
	p.registerPrefixParser(token.MINUS, p.parsePrefixExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := ast.NewProgram()

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStmt()
		if stmt != nil {
			program.AppendStmt(stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStmt() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStmt()
	case token.RETURN:
		return p.parseReturnStmt()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStmt() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: skipping the expression
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	msg := "expected next token to be %s, got %s instead"
	p.addError(fmt.Sprintf(msg, t, p.peekToken.Type))
}

func (p *Parser) parseReturnStmt() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) registerPrefixParser(tt token.TokenType, pp prefixParser) {
	p.prefixParsers[tt] = pp
}

func (p *Parser) registerInfixParser(tt token.TokenType, ip infixParser) {
	p.infixParsers[tt] = ip
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token:      p.curToken,
		Expression: p.parseExpression(PRECEDENCE_LOWEST),
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(op OperatorPrecedence) ast.Expression {
	parser, ok := p.prefixParsers[p.curToken.Type]
	if !ok {
		p.parserNotFound(p.curToken.Type)
		return nil
	}

	return parser()
}

func (p *Parser) addError(errMsg string) {
	p.errors = append(p.errors, errMsg)
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()

	exp.Right = p.parseExpression(PRECEDENCE_PREFIX)
	return exp
}
