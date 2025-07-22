package eval

import (
	"monkey/ast"
	"monkey/object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.PrefixExpression:
		return evalPrefixExp(node, env)
	case *ast.InfixExpression:
		return evalInfixExp(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.BlockStatement:
		return evalBlockStmt(node, env)
	case *ast.ReturnStatement:
		return evalReturn(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return object.AsBool(node.Value)
	case *ast.LetStatement:
		return evalLetStmt(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	default:
		return nil
	}
}

func evalIdentifier(ident *ast.Identifier, env *object.Environment) object.Object {
	obj, ok := env.Get(ident.Value)

	if !ok {
		return object.FormatError("identifier not found: %s", ident.Value)
	}

	return obj
}
func evalLetStmt(ls *ast.LetStatement, env *object.Environment) object.Object {
	val := Eval(ls.Value, env)
	if object.IsError(val) {
		return val
	}
	return env.Set(ls.Name.Value, val)
}

func evalReturn(rv *ast.ReturnStatement, env *object.Environment) object.Object {
	val := Eval(rv.ReturnValue, env)
	if object.IsError(val) {
		return val
	}
	return &object.ReturnValue{Value: val}
}

func evalProgram(p *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range p.Statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalPrefixExp(node *ast.PrefixExpression, env *object.Environment) object.Object {
	val := Eval(node.Right, env)
	if object.IsError(val) {
		return val
	}

	switch node.Operator {
	case "!":
		return evalBang(val)
	case "-":
		return evalMinus(val)
	default:
		return object.FormatError(
			"unknown operator: %s%s", node.Operator, val.Type())
	}
}

func evalBang(right object.Object) object.Object {
	switch right {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NULL:
		return object.TRUE
	default:
		return object.FALSE
	}
}

func evalMinus(right object.Object) object.Object {
	if right.Type() != object.OBJ_INTEGER {
		return object.FormatError("unknown operator: -%s",
			right.Type())
	}
	integer := right.(*object.Integer)
	return object.AsInt(-integer.Value)
}

func evalInfixExp(node *ast.InfixExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	right := Eval(node.Right, env)

	if object.IsError(left) {
		return left
	}

	if object.IsError(right) {
		return right
	}

	switch {
	case left.Type() == object.OBJ_INTEGER && right.Type() == object.OBJ_INTEGER:
		return evalIntegerInfixExp(left, right, node.Operator)
	case left.Type() != right.Type():
		return object.FormatError("type mismatch: %s %s %s",
			left.Type(), node.Operator, right.Type())
	case node.Operator == "==":
		return object.AsBool(left == right)
	case node.Operator == "!=":
		return object.AsBool(left != right)
	default:
		return object.FormatError("unknown operator: %s %s %s",
			left.Type(), node.Operator, right.Type())
	}
}

func evalIntegerInfixExp(left, right object.Object, operator string) object.Object {
	leftInt := left.(*object.Integer).Value
	rightInt := right.(*object.Integer).Value

	switch operator {
	case "+":
		return object.AsInt(leftInt + rightInt)
	case "-":
		return object.AsInt(leftInt - rightInt)
	case "*":
		return object.AsInt(leftInt * rightInt)
	case "/":
		return object.AsInt(leftInt / rightInt)
	case ">":
		return object.AsBool(leftInt > rightInt)
	case "<":
		return object.AsBool(leftInt < rightInt)
	case "==":
		return object.AsBool(leftInt == rightInt)
	case "!=":
		return object.AsBool(leftInt != rightInt)
	default:
		return object.FormatError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ifExp *ast.IfExpression, env *object.Environment) object.Object {
	cond := Eval(ifExp.Condition, env)

	if object.IsTruthy(cond) {
		return Eval(ifExp.Consequence, env)
	}
	if ifExp.Alternative != nil {
		return Eval(ifExp.Alternative, env)
	}
	return object.NULL
}

func evalBlockStmt(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		if result == nil {
			continue
		}

		switch result.Type() {
		case object.OBJ_RETURN_VALUE:
			fallthrough
		case object.OBJ_ERROR:
			return result
		}
	}

	return result
}
