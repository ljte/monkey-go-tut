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
	p.registerPrefixParser(token.TRUE, p.parseBoolean)
	p.registerPrefixParser(token.FALSE, p.parseBoolean)
	p.registerPrefixParser(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefixParser(token.IF, p.parseIfExpression)
	p.registerPrefixParser(token.FUNCTION, p.parseFunctionLiteral)

	p.registerInfixParser(token.PLUS, p.parseInfixExpression)
	p.registerInfixParser(token.MINUS, p.parseInfixExpression)
	p.registerInfixParser(token.SLASH, p.parseInfixExpression)
	p.registerInfixParser(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixParser(token.EQUALS, p.parseInfixExpression)
	p.registerInfixParser(token.NOT_EQUALS, p.parseInfixExpression)
	p.registerInfixParser(token.LT, p.parseInfixExpression)
	p.registerInfixParser(token.GT, p.parseInfixExpression)
	p.registerInfixParser(token.LPAREN, p.parseCallExpression)

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
	p.nextToken()

	stmt.Value = p.parseExpression(PRECEDENCE_LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
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

	stmt.ReturnValue = p.parseExpression(PRECEDENCE_LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
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

func (p *Parser) parseExpression(precedence OperatorPrecedence) ast.Expression {
	parser, ok := p.prefixParsers[p.curToken.Type]
	if !ok {
		p.prefixParserNotFound(p.curToken.Type)
		return nil
	}

	left := parser()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peerPrecedence() {
		infix := p.infixParsers[p.peekToken.Type]
		if infix == nil {
			return left
		}

		p.nextToken()
		left = infix(left)
	}
	return left
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

func (p *Parser) peerPrecedence() OperatorPrecedence {
	return DerivePrecedence(p.peekToken.Type)
}

func (p *Parser) curPrecedence() OperatorPrecedence {
	return DerivePrecedence(p.curToken.Type)
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precedence)
	return exp
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(PRECEDENCE_LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	exp.Condition = p.parseExpression(PRECEDENCE_LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LCURLY) {
		return nil
	}

	exp.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LCURLY) {
			return nil
		}

		exp.Alternative = p.parseBlockStatement()
	}
	return exp
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RCURLY) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStmt()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LCURLY) {
		return nil
	}

	lit.Body = p.parseBlockStatement()
	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}
	p.nextToken()

	identifiers = append(identifiers, &ast.Identifier{
		Token: p.curToken, Value: p.curToken.Literal,
	})

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		identifiers = append(identifiers, &ast.Identifier{
			Token: p.curToken, Value: p.curToken.Literal,
		})
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	return &ast.CallExpression{
		Token: p.curToken,
		Func:  function,
		Args:  p.parseCallArguments(),
	}
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(PRECEDENCE_LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(PRECEDENCE_LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}
