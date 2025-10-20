package main

import (
	"fmt"
	"unicode"
)

type TokenType int

const (
	// Literals
	NUMBER TokenType = iota
	STRING
	IDENTIFIER
	BOOLEAN

	// F1-themed keywords
	GRID          // var (starting grid position)
	PACE          // func (racing pace/function)
	CIRCUIT       // if (racing circuit/conditional)
	ELSE_CIRCUIT  // else
	LOOP          // for (racing loop)
	WHILE_RACING  // while
	RETURN_PIT    // return (return to pit)
	BREAK_FLAG    // break (yellow flag)
	CONTINUE_RACE // continue

	// Data structures
	FORMATION // array (formation lap)
	GARAGE    // struct/object

	// Operators
	ASSIGN   // =
	PLUS     // +
	MINUS    // -
	MULTIPLY // *
	DIVIDE   // /
	MODULO   // %

	// Comparison
	EQUAL         // ==
	NOT_EQUAL     // !=
	LESS          // <
	LESS_EQUAL    // <=
	GREATER       // >
	GREATER_EQUAL // >=

	// Logical
	AND // &&
	OR  // ||
	NOT // !

	// Delimiters
	SEMICOLON // ;
	COMMA     // ,
	DOT       // .

	// Brackets
	LPAREN   // (
	RPAREN   // )
	LBRACE   // {
	RBRACE   // }
	LBRACKET // [
	RBRACKET // ]

	// Special
	NEWLINE
	EOF
	ILLEGAL
)

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

type Lexer struct {
	input    string
	position int
	line     int
	column   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		line:   1,
		column: 1,
	}
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	if l.position >= len(l.input) {
		return Token{Type: EOF, Line: l.line, Column: l.column}
	}

	ch := l.input[l.position]

	// Handle comments
	if ch == '/' && l.peek() == '/' {
		l.skipComment()
		return l.NextToken()
	}

	switch ch {
	case '\n':
		token := Token{Type: NEWLINE, Value: string(ch), Line: l.line, Column: l.column}
		l.advance()
		l.line++
		l.column = 1
		return token
	case '=':
		if l.peek() == '=' {
			l.advance()
			l.advance()
			return Token{Type: EQUAL, Value: "==", Line: l.line, Column: l.column - 2}
		}
		return l.singleCharToken(ASSIGN)
	case '+':
		return l.singleCharToken(PLUS)
	case '-':
		return l.singleCharToken(MINUS)
	case '*':
		return l.singleCharToken(MULTIPLY)
	case '/':
		return l.singleCharToken(DIVIDE)
	case '%':
		return l.singleCharToken(MODULO)
	case '!':
		if l.peek() == '=' {
			l.advance()
			l.advance()
			return Token{Type: NOT_EQUAL, Value: "!=", Line: l.line, Column: l.column - 2}
		}
		return l.singleCharToken(NOT)
	case '<':
		if l.peek() == '=' {
			l.advance()
			l.advance()
			return Token{Type: LESS_EQUAL, Value: "<=", Line: l.line, Column: l.column - 2}
		}
		return l.singleCharToken(LESS)
	case '>':
		if l.peek() == '=' {
			l.advance()
			l.advance()
			return Token{Type: GREATER_EQUAL, Value: ">=", Line: l.line, Column: l.column - 2}
		}
		return l.singleCharToken(GREATER)
	case '&':
		if l.peek() == '&' {
			l.advance()
			l.advance()
			return Token{Type: AND, Value: "&&", Line: l.line, Column: l.column - 2}
		}
		return Token{Type: ILLEGAL, Value: string(ch), Line: l.line, Column: l.column}
	case '|':
		if l.peek() == '|' {
			l.advance()
			l.advance()
			return Token{Type: OR, Value: "||", Line: l.line, Column: l.column - 2}
		}
		return Token{Type: ILLEGAL, Value: string(ch), Line: l.line, Column: l.column}
	case ';':
		return l.singleCharToken(SEMICOLON)
	case ',':
		return l.singleCharToken(COMMA)
	case '.':
		return l.singleCharToken(DOT)
	case '(':
		return l.singleCharToken(LPAREN)
	case ')':
		return l.singleCharToken(RPAREN)
	case '{':
		return l.singleCharToken(LBRACE)
	case '}':
		return l.singleCharToken(RBRACE)
	case '[':
		return l.singleCharToken(LBRACKET)
	case ']':
		return l.singleCharToken(RBRACKET)
	case '"':
		return l.readString()
	default:
		if unicode.IsDigit(rune(ch)) {
			return l.readNumber()
		}
		if unicode.IsLetter(rune(ch)) || ch == '_' {
			return l.readIdentifier()
		}
		return Token{Type: ILLEGAL, Value: string(ch), Line: l.line, Column: l.column}
	}
}

func (l *Lexer) singleCharToken(tokenType TokenType) Token {
	ch := l.input[l.position]
	token := Token{Type: tokenType, Value: string(ch), Line: l.line, Column: l.column}
	l.advance()
	return token
}

func (l *Lexer) advance() {
	l.position++
	l.column++
}

func (l *Lexer) peek() byte {
	if l.position+1 >= len(l.input) {
		return 0
	}
	return l.input[l.position+1]
}

func (l *Lexer) skipWhitespace() {
	for l.position < len(l.input) {
		ch := l.input[l.position]
		if ch == ' ' || ch == '\t' || ch == '\r' {
			l.advance()
		} else {
			break
		}
	}
}

func (l *Lexer) skipComment() {
	for l.position < len(l.input) && l.input[l.position] != '\n' {
		l.advance()
	}
}

func (l *Lexer) readString() Token {
	start := l.position
	startCol := l.column
	l.advance() // skip opening quote

	for l.position < len(l.input) && l.input[l.position] != '"' {
		l.advance()
	}

	if l.position >= len(l.input) {
		return Token{Type: ILLEGAL, Value: "unterminated string", Line: l.line, Column: startCol}
	}

	value := l.input[start+1 : l.position]
	l.advance() // skip closing quote

	return Token{Type: STRING, Value: value, Line: l.line, Column: startCol}
}

func (l *Lexer) readNumber() Token {
	start := l.position
	startCol := l.column

	for l.position < len(l.input) && (unicode.IsDigit(rune(l.input[l.position])) || l.input[l.position] == '.') {
		l.advance()
	}

	value := l.input[start:l.position]
	return Token{Type: NUMBER, Value: value, Line: l.line, Column: startCol}
}

func (l *Lexer) readIdentifier() Token {
	start := l.position
	startCol := l.column

	for l.position < len(l.input) && (unicode.IsLetter(rune(l.input[l.position])) || unicode.IsDigit(rune(l.input[l.position])) || l.input[l.position] == '_') {
		l.advance()
	}

	value := l.input[start:l.position]
	tokenType := l.getKeywordType(value)

	return Token{Type: tokenType, Value: value, Line: l.line, Column: startCol}
}

func (l *Lexer) getKeywordType(identifier string) TokenType {
	keywords := map[string]TokenType{
		"grid":          GRID,
		"pace":          PACE,
		"circuit":       CIRCUIT,
		"else_circuit":  ELSE_CIRCUIT,
		"loop":          LOOP,
		"while_racing":  WHILE_RACING,
		"return_pit":    RETURN_PIT,
		"break_flag":    BREAK_FLAG,
		"continue_race": CONTINUE_RACE,
		"formation":     FORMATION,
		"garage":        GARAGE,
		"true":          BOOLEAN,
		"false":         BOOLEAN,
	}

	if tokenType, exists := keywords[identifier]; exists {
		return tokenType
	}

	return IDENTIFIER
}

func (t TokenType) String() string {
	names := map[TokenType]string{
		NUMBER: "NUMBER", STRING: "STRING", IDENTIFIER: "IDENTIFIER", BOOLEAN: "BOOLEAN",
		GRID: "GRID", PACE: "PACE", CIRCUIT: "CIRCUIT", ELSE_CIRCUIT: "ELSE_CIRCUIT",
		LOOP: "LOOP", WHILE_RACING: "WHILE_RACING", RETURN_PIT: "RETURN_PIT",
		BREAK_FLAG: "BREAK_FLAG", CONTINUE_RACE: "CONTINUE_RACE",
		FORMATION: "FORMATION", GARAGE: "GARAGE",
		ASSIGN: "ASSIGN", PLUS: "PLUS", MINUS: "MINUS", MULTIPLY: "MULTIPLY", DIVIDE: "DIVIDE", MODULO: "MODULO",
		EQUAL: "EQUAL", NOT_EQUAL: "NOT_EQUAL", LESS: "LESS", LESS_EQUAL: "LESS_EQUAL",
		GREATER: "GREATER", GREATER_EQUAL: "GREATER_EQUAL",
		AND: "AND", OR: "OR", NOT: "NOT",
		SEMICOLON: "SEMICOLON", COMMA: "COMMA", DOT: "DOT",
		LPAREN: "LPAREN", RPAREN: "RPAREN", LBRACE: "LBRACE", RBRACE: "RBRACE",
		LBRACKET: "LBRACKET", RBRACKET: "RBRACKET",
		NEWLINE: "NEWLINE", EOF: "EOF", ILLEGAL: "ILLEGAL",
	}

	if name, exists := names[t]; exists {
		return name
	}
	return fmt.Sprintf("UNKNOWN(%d)", int(t))
}
