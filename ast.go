package main

import "fmt"

// AST Node interface
type Node interface {
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Statements
type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	result := ""
	for _, stmt := range p.Statements {
		result += stmt.String()
	}
	return result
}

type GridStatement struct { // var declaration
	Name  *Identifier
	Value Expression
}

func (gs *GridStatement) statementNode() {}
func (gs *GridStatement) String() string {
	return fmt.Sprintf("grid %s = %s;", gs.Name.String(), gs.Value.String())
}

type PaceStatement struct { // function declaration
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (ps *PaceStatement) statementNode() {}
func (ps *PaceStatement) String() string {
	params := ""
	for i, p := range ps.Parameters {
		if i > 0 {
			params += ", "
		}
		params += p.String()
	}
	return fmt.Sprintf("pace %s(%s) %s", ps.Name.String(), params, ps.Body.String())
}

type CircuitStatement struct { // if statement
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (cs *CircuitStatement) statementNode() {}
func (cs *CircuitStatement) String() string {
	result := fmt.Sprintf("circuit (%s) %s", cs.Condition.String(), cs.Consequence.String())
	if cs.Alternative != nil {
		result += fmt.Sprintf(" else_circuit %s", cs.Alternative.String())
	}
	return result
}

type LoopStatement struct { // for loop
	Init      Statement
	Condition Expression
	Update    Statement
	Body      *BlockStatement
}

func (ls *LoopStatement) statementNode() {}
func (ls *LoopStatement) String() string {
	return fmt.Sprintf("loop (%s; %s; %s) %s",
		ls.Init.String(), ls.Condition.String(), ls.Update.String(), ls.Body.String())
}

type WhileRacingStatement struct { // while loop
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileRacingStatement) statementNode() {}
func (ws *WhileRacingStatement) String() string {
	return fmt.Sprintf("while_racing (%s) %s", ws.Condition.String(), ws.Body.String())
}

type ReturnPitStatement struct { // return statement
	Value Expression
}

func (rs *ReturnPitStatement) statementNode() {}
func (rs *ReturnPitStatement) String() string {
	if rs.Value != nil {
		return fmt.Sprintf("return_pit %s;", rs.Value.String())
	}
	return "return_pit;"
}

type BreakFlagStatement struct{} // break statement

func (bs *BreakFlagStatement) statementNode() {}
func (bs *BreakFlagStatement) String() string {
	return "break_flag;"
}

type ContinueRaceStatement struct{} // continue statement

func (cs *ContinueRaceStatement) statementNode() {}
func (cs *ContinueRaceStatement) String() string {
	return "continue_race;"
}

type ExpressionStatement struct {
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string {
	return es.Expression.String() + ";"
}

type BlockStatement struct {
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) String() string {
	result := "{"
	for _, stmt := range bs.Statements {
		result += stmt.String()
	}
	result += "}"
	return result
}

// Expressions
type Identifier struct {
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string {
	return i.Value
}

type NumberLiteral struct {
	Value float64
}

func (nl *NumberLiteral) expressionNode() {}
func (nl *NumberLiteral) String() string {
	return fmt.Sprintf("%g", nl.Value)
}

type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string {
	return fmt.Sprintf("\"%s\"", sl.Value)
}

type BooleanLiteral struct {
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) String() string {
	if bl.Value {
		return "true"
	}
	return "false"
}

type FormationLiteral struct { // array literal
	Elements []Expression
}

func (fl *FormationLiteral) expressionNode() {}
func (fl *FormationLiteral) String() string {
	result := "["
	for i, elem := range fl.Elements {
		if i > 0 {
			result += ", "
		}
		result += elem.String()
	}
	result += "]"
	return result
}

type IndexExpression struct {
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) String() string {
	return fmt.Sprintf("(%s[%s])", ie.Left.String(), ie.Index.String())
}

type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator, ie.Right.String())
}

type PrefixExpression struct {
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right.String())
}

type CallExpression struct {
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) String() string {
	args := ""
	for i, arg := range ce.Arguments {
		if i > 0 {
			args += ", "
		}
		args += arg.String()
	}
	return fmt.Sprintf("%s(%s)", ce.Function.String(), args)
}

type AssignmentExpression struct {
	Name  *Identifier
	Value Expression
}

func (ae *AssignmentExpression) expressionNode() {}
func (ae *AssignmentExpression) String() string {
	return fmt.Sprintf("%s = %s", ae.Name.String(), ae.Value.String())
}
