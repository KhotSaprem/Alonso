package main

import (
	"fmt"
	"strings"
)

type ObjectType string

const (
	NUMBER_OBJ   = "NUMBER"
	STRING_OBJ   = "STRING"
	BOOLEAN_OBJ  = "BOOLEAN"
	NULL_OBJ     = "NULL"
	RETURN_OBJ   = "RETURN_VALUE"
	ERROR_OBJ    = "ERROR"
	FUNCTION_OBJ = "FUNCTION"
	BUILTIN_OBJ  = "BUILTIN"
	ARRAY_OBJ    = "ARRAY"
	BREAK_OBJ    = "BREAK"
	CONTINUE_OBJ = "CONTINUE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Number struct {
	Value float64
}

func (n *Number) Type() ObjectType { return NUMBER_OBJ }
func (n *Number) Inspect() string  { return fmt.Sprintf("%g", n.Value) }

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Function struct {
	Parameters []*Identifier
	Body       *BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	return fmt.Sprintf("pace(%s) {\n%s\n}", strings.Join(params, ", "), f.Body.String())
}

type Builtin struct {
	Fn func(args ...Object) Object
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

type Break struct{}

func (b *Break) Type() ObjectType { return BREAK_OBJ }
func (b *Break) Inspect() string  { return "break_flag" }

type Continue struct{}

func (c *Continue) Type() ObjectType { return CONTINUE_OBJ }
func (c *Continue) Inspect() string  { return "continue_race" }

// Environment for variable scoping
type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	value, ok := e.store[name]
	if !ok && e.outer != nil {
		value, ok = e.outer.Get(name)
	}
	return value, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
