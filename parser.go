package main

import (
	"fmt"
	"strconv"
)

type Parser struct {
	lexer        *Lexer
	currentToken Token
	peekToken    Token
	errors       []string
}

type PrecedenceLevel int

const (
	_ PrecedenceLevel = iota
	LOWEST
	LOGICAL     // && ||
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var precedences = map[TokenType]PrecedenceLevel{
	AND:           LOGICAL,
	OR:            LOGICAL,
	EQUAL:         EQUALS,
	NOT_EQUAL:     EQUALS,
	LESS:          LESSGREATER,
	GREATER:       LESSGREATER,
	LESS_EQUAL:    LESSGREATER,
	GREATER_EQUAL: LESSGREATER,
	PLUS:          SUM,
	MINUS:         SUM,
	DIVIDE:        PRODUCT,
	MULTIPLY:      PRODUCT,
	MODULO:        PRODUCT,
	ASSIGN:        EQUALS,
	LPAREN:        CALL,
	LBRACKET:      INDEX,
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	// Read two tokens, so currentToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for p.currentToken.Type != EOF {
		// Skip newlines
		if p.currentToken.Type == NEWLINE {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.currentToken.Type {
	case GRID:
		return p.parseGridStatement()
	case PACE:
		return p.parsePaceStatement()
	case CIRCUIT:
		return p.parseCircuitStatement()
	case LOOP:
		return p.parseLoopStatement()
	case WHILE_RACING:
		return p.parseWhileRacingStatement()
	case RETURN_PIT:
		return p.parseReturnPitStatement()
	case BREAK_FLAG:
		return p.parseBreakFlagStatement()
	case CONTINUE_RACE:
		return p.parseContinueRaceStatement()
	case LBRACE:
		return p.parseBlockStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseGridStatement() *GridStatement {
	stmt := &GridStatement{}

	if !p.expectPeek(IDENTIFIER) {
		return nil
	}

	stmt.Name = &Identifier{Value: p.currentToken.Value}

	if !p.expectPeek(ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekToken.Type == SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parsePaceStatement() *PaceStatement {
	stmt := &PaceStatement{}

	if !p.expectPeek(IDENTIFIER) {
		return nil
	}

	stmt.Name = &Identifier{Value: p.currentToken.Value}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseFunctionParameters() []*Identifier {
	identifiers := []*Identifier{}

	if p.peekToken.Type == RPAREN {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &Identifier{Value: p.currentToken.Value}
	identifiers = append(identifiers, ident)

	for p.peekToken.Type == COMMA {
		p.nextToken()
		p.nextToken()
		ident := &Identifier{Value: p.currentToken.Value}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCircuitStatement() *CircuitStatement {
	stmt := &CircuitStatement{}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(RPAREN) {
		return nil
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	stmt.Consequence = p.parseBlockStatement()

	if p.peekToken.Type == ELSE_CIRCUIT {
		p.nextToken()

		if !p.expectPeek(LBRACE) {
			return nil
		}

		stmt.Alternative = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parseLoopStatement() *LoopStatement {
	stmt := &LoopStatement{}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Init = p.parseStatement()

	if p.currentToken.Type == SEMICOLON {
		p.nextToken()
	}

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(SEMICOLON) {
		return nil
	}

	p.nextToken()
	stmt.Update = p.parseExpressionStatement()

	if !p.expectPeek(RPAREN) {
		return nil
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseWhileRacingStatement() *WhileRacingStatement {
	stmt := &WhileRacingStatement{}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(RPAREN) {
		return nil
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseReturnPitStatement() *ReturnPitStatement {
	stmt := &ReturnPitStatement{}

	if p.peekToken.Type != SEMICOLON && p.peekToken.Type != NEWLINE {
		p.nextToken()
		stmt.Value = p.parseExpression(LOWEST)
	}

	if p.peekToken.Type == SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBreakFlagStatement() *BreakFlagStatement {
	if p.peekToken.Type == SEMICOLON {
		p.nextToken()
	}
	return &BreakFlagStatement{}
}

func (p *Parser) parseContinueRaceStatement() *ContinueRaceStatement {
	if p.peekToken.Type == SEMICOLON {
		p.nextToken()
	}
	return &ContinueRaceStatement{}
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{}
	block.Statements = []Statement{}

	p.nextToken()

	for p.currentToken.Type != RBRACE && p.currentToken.Type != EOF {
		if p.currentToken.Type == NEWLINE {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekToken.Type == SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence PrecedenceLevel) Expression {
	var leftExp Expression

	switch p.currentToken.Type {
	case IDENTIFIER:
		leftExp = p.parseIdentifier()
	case NUMBER:
		leftExp = p.parseNumberLiteral()
	case STRING:
		leftExp = p.parseStringLiteral()
	case BOOLEAN:
		leftExp = p.parseBooleanLiteral()
	case LBRACKET:
		leftExp = p.parseFormationLiteral()
	case MINUS, NOT:
		leftExp = p.parsePrefixExpression()
	case LPAREN:
		leftExp = p.parseGroupedExpression()
	default:
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}

	for p.peekToken.Type != SEMICOLON && precedence < p.peekPrecedence() {
		switch p.peekToken.Type {
		case PLUS, MINUS, DIVIDE, MULTIPLY, MODULO, EQUAL, NOT_EQUAL, LESS, GREATER, LESS_EQUAL, GREATER_EQUAL, AND, OR:
			p.nextToken()
			leftExp = p.parseInfixExpression(leftExp)
		case LPAREN:
			p.nextToken()
			leftExp = p.parseCallExpression(leftExp)
		case LBRACKET:
			p.nextToken()
			leftExp = p.parseIndexExpression(leftExp)
		case ASSIGN:
			p.nextToken()
			leftExp = p.parseAssignmentExpression(leftExp)
		default:
			return leftExp
		}
	}

	return leftExp
}

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Value: p.currentToken.Value}
}

func (p *Parser) parseNumberLiteral() Expression {
	lit := &NumberLiteral{}

	value, err := strconv.ParseFloat(p.currentToken.Value, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as number", p.currentToken.Value)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Value: p.currentToken.Value}
}

func (p *Parser) parseBooleanLiteral() Expression {
	return &BooleanLiteral{Value: p.currentToken.Value == "true"}
}

func (p *Parser) parseFormationLiteral() Expression {
	lit := &FormationLiteral{}
	lit.Elements = p.parseExpressionList(RBRACKET)
	return lit
}

func (p *Parser) parseExpressionList(end TokenType) []Expression {
	args := []Expression{}

	if p.peekToken.Type == end {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekToken.Type == COMMA {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return args
}

func (p *Parser) parsePrefixExpression() Expression {
	expression := &PrefixExpression{
		Operator: p.currentToken.Value,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Left:     left,
		Operator: p.currentToken.Value,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseCallExpression(fn Expression) Expression {
	exp := &CallExpression{Function: fn}
	exp.Arguments = p.parseExpressionList(RPAREN)
	return exp
}

func (p *Parser) parseIndexExpression(left Expression) Expression {
	exp := &IndexExpression{Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseAssignmentExpression(left Expression) Expression {
	ident, ok := left.(*Identifier)
	if !ok {
		p.errors = append(p.errors, "invalid assignment target")
		return nil
	}

	exp := &AssignmentExpression{Name: ident}

	p.nextToken()
	exp.Value = p.parseExpression(LOWEST)

	return exp
}

func (p *Parser) curPrecedence() PrecedenceLevel {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() PrecedenceLevel {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}
