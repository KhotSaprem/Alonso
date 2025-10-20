package main

import (
	"fmt"
)

var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

type Interpreter struct {
	env *Environment
}

func NewInterpreter() *Interpreter {
	env := NewEnvironment()

	// Add built-in functions
	builtins := map[string]*Builtin{
		"telemetry": { // print function
			Fn: func(args ...Object) Object {
				for i, arg := range args {
					if i > 0 {
						fmt.Print(" ")
					}
					fmt.Print(arg.Inspect())
				}
				fmt.Println()
				return NULL
			},
		},
		"length": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				switch arg := args[0].(type) {
				case *Array:
					return &Number{Value: float64(len(arg.Elements))}
				case *String:
					return &Number{Value: float64(len(arg.Value))}
				default:
					return newError("argument to `length` not supported, got %T", arg)
				}
			},
		},
		"push": {
			Fn: func(args ...Object) Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}

				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to `push` must be ARRAY, got %T", args[0])
				}

				arr := args[0].(*Array)
				length := len(arr.Elements)

				newElements := make([]Object, length+1)
				copy(newElements, arr.Elements)
				newElements[length] = args[1]

				return &Array{Elements: newElements}
			},
		},
	}

	for name, builtin := range builtins {
		env.Set(name, builtin)
	}

	return &Interpreter{env: env}
}

func (i *Interpreter) Execute(input string) error {
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		for _, err := range parser.Errors() {
			fmt.Printf("Parser error: %s\n", err)
		}
		return fmt.Errorf("parsing failed")
	}

	result := i.Eval(program, i.env)
	if result != nil && result.Type() == ERROR_OBJ {
		return fmt.Errorf(result.Inspect())
	}

	return nil
}

func (i *Interpreter) Eval(node Node, env *Environment) Object {
	switch node := node.(type) {

	// Statements
	case *Program:
		return i.evalProgram(node.Statements, env)

	case *GridStatement:
		val := i.Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return val

	case *PaceStatement:
		fn := &Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
		env.Set(node.Name.Value, fn)
		return fn

	case *CircuitStatement:
		return i.evalCircuitStatement(node, env)

	case *LoopStatement:
		return i.evalLoopStatement(node, env)

	case *WhileRacingStatement:
		return i.evalWhileRacingStatement(node, env)

	case *ReturnPitStatement:
		var val Object = NULL
		if node.Value != nil {
			evaluated := i.Eval(node.Value, env)
			if isError(evaluated) {
				return evaluated
			}
			val = evaluated
		}
		return &ReturnValue{Value: val}

	case *BreakFlagStatement:
		return &Break{}

	case *ContinueRaceStatement:
		return &Continue{}

	case *BlockStatement:
		return i.evalBlockStatement(node, env)

	case *ExpressionStatement:
		return i.Eval(node.Expression, env)

	// Expressions
	case *NumberLiteral:
		return &Number{Value: node.Value}

	case *StringLiteral:
		return &String{Value: node.Value}

	case *BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)

	case *FormationLiteral:
		elements := i.evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &Array{Elements: elements}

	case *IndexExpression:
		left := i.Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := i.Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return i.evalIndexExpression(left, index)

	case *Identifier:
		return i.evalIdentifier(node, env)

	case *PrefixExpression:
		right := i.Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return i.evalPrefixExpression(node.Operator, right)

	case *InfixExpression:
		left := i.Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := i.Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return i.evalInfixExpression(node.Operator, left, right)

	case *CallExpression:
		function := i.Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := i.evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return i.applyFunction(function, args)

	case *AssignmentExpression:
		val := i.Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return val

	default:
		return newError("unknown node type: %T", node)
	}
}

func (i *Interpreter) evalProgram(stmts []Statement, env *Environment) Object {
	var result Object

	for _, statement := range stmts {
		result = i.Eval(statement, env)

		switch result := result.(type) {
		case *ReturnValue:
			return result.Value
		case *Error:
			return result
		}
	}

	return result
}

func (i *Interpreter) evalBlockStatement(block *BlockStatement, env *Environment) Object {
	var result Object

	for _, statement := range block.Statements {
		result = i.Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == RETURN_OBJ || rt == ERROR_OBJ || rt == BREAK_OBJ || rt == CONTINUE_OBJ {
				return result
			}
		}
	}

	return result
}

func (i *Interpreter) evalCircuitStatement(node *CircuitStatement, env *Environment) Object {
	condition := i.Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return i.Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return i.Eval(node.Alternative, env)
	} else {
		return NULL
	}
}

func (i *Interpreter) evalLoopStatement(node *LoopStatement, env *Environment) Object {
	// Create new scope for loop
	loopEnv := NewEnclosedEnvironment(env)

	// Initialize
	if node.Init != nil {
		result := i.Eval(node.Init, loopEnv)
		if isError(result) {
			return result
		}
	}

	var result Object = NULL

	for {
		// Check condition
		if node.Condition != nil {
			condition := i.Eval(node.Condition, loopEnv)
			if isError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		// Execute body
		result = i.Eval(node.Body, loopEnv)
		if result != nil {
			switch result.Type() {
			case RETURN_OBJ, ERROR_OBJ:
				return result
			case BREAK_OBJ:
				return NULL
			case CONTINUE_OBJ:
				// Continue to update
			}
		}

		// Update
		if node.Update != nil {
			updateResult := i.Eval(node.Update, loopEnv)
			if isError(updateResult) {
				return updateResult
			}
		}
	}

	return result
}

func (i *Interpreter) evalWhileRacingStatement(node *WhileRacingStatement, env *Environment) Object {
	var result Object = NULL

	for {
		condition := i.Eval(node.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = i.Eval(node.Body, env)
		if result != nil {
			switch result.Type() {
			case RETURN_OBJ, ERROR_OBJ:
				return result
			case BREAK_OBJ:
				return NULL
			case CONTINUE_OBJ:
				continue
			}
		}
	}

	return result
}

func (i *Interpreter) evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "!":
		return i.evalBangOperatorExpression(right)
	case "-":
		return i.evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%T", operator, right)
	}
}

func (i *Interpreter) evalBangOperatorExpression(right Object) Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func (i *Interpreter) evalMinusPrefixOperatorExpression(right Object) Object {
	if right.Type() != NUMBER_OBJ {
		return newError("unknown operator: -%T", right)
	}

	value := right.(*Number).Value
	return &Number{Value: -value}
}

func (i *Interpreter) evalInfixExpression(operator string, left, right Object) Object {
	switch {
	case left.Type() == NUMBER_OBJ && right.Type() == NUMBER_OBJ:
		return i.evalNumberInfixExpression(operator, left, right)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return i.evalStringInfixExpression(operator, left, right)
	case operator == "&&":
		return nativeBoolToBooleanObject(isTruthy(left) && isTruthy(right))
	case operator == "||":
		return nativeBoolToBooleanObject(isTruthy(left) || isTruthy(right))
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %T %s %T", left, operator, right)
	default:
		return newError("unknown operator: %T %s %T", left, operator, right)
	}
}

func (i *Interpreter) evalNumberInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*Number).Value
	rightVal := right.(*Number).Value

	switch operator {
	case "+":
		return &Number{Value: leftVal + rightVal}
	case "-":
		return &Number{Value: leftVal - rightVal}
	case "*":
		return &Number{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &Number{Value: leftVal / rightVal}
	case "%":
		return &Number{Value: float64(int(leftVal) % int(rightVal))}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s", operator)
	}
}

func (i *Interpreter) evalStringInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*String).Value
	rightVal := right.(*String).Value

	switch operator {
	case "+":
		return &String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %T %s %T", left, operator, right)
	}
}

func (i *Interpreter) evalIndexExpression(left, index Object) Object {
	switch {
	case left.Type() == ARRAY_OBJ && index.Type() == NUMBER_OBJ:
		return i.evalArrayIndexExpression(left, index)
	default:
		return newError("index operator not supported: %T", left)
	}
}

func (i *Interpreter) evalArrayIndexExpression(array, index Object) Object {
	arrayObject := array.(*Array)
	idx := int(index.(*Number).Value)
	max := len(arrayObject.Elements) - 1

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func (i *Interpreter) evalIdentifier(node *Identifier, env *Environment) Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: " + node.Value)
	}

	return val
}

func (i *Interpreter) evalExpressions(exps []Expression, env *Environment) []Object {
	result := []Object{}

	for _, e := range exps {
		evaluated := i.Eval(e, env)
		if isError(evaluated) {
			return []Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func (i *Interpreter) applyFunction(fn Object, args []Object) Object {
	switch fn := fn.(type) {
	case *Function:
		extendedEnv := i.extendFunctionEnv(fn, args)
		evaluated := i.Eval(fn.Body, extendedEnv)
		return i.unwrapReturnValue(evaluated)
	case *Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %T", fn)
	}
}

func (i *Interpreter) extendFunctionEnv(fn *Function, args []Object) *Environment {
	env := NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		if paramIdx < len(args) {
			env.Set(param.Value, args[paramIdx])
		}
	}

	return env
}

func (i *Interpreter) unwrapReturnValue(obj Object) Object {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func nativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
