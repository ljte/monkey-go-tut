package eval

import (
	"monkey/ast"
	"monkey/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		return evalPrefixExp(node)
	case *ast.InfixExpression:
		return evalInfixExp(node)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.BlockStatement:
		return evalBlockStmt(node)
	case *ast.ReturnStatement:
		return evalReturn(node)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return object.AsBool(node.Value)
	default:
		return nil
	}
}

func evalReturn(rv *ast.ReturnStatement) object.Object {
	val := Eval(rv.ReturnValue)
	if object.IsError(val) {
		return val
	}
	return &object.ReturnValue{Value: val}
}

func evalProgram(p *ast.Program) object.Object {
	var result object.Object

	for _, stmt := range p.Statements {
		result = Eval(stmt)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalPrefixExp(node *ast.PrefixExpression) object.Object {
	val := Eval(node.Right)
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

func evalInfixExp(node *ast.InfixExpression) object.Object {
	left := Eval(node.Left)
	right := Eval(node.Right)

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

func evalIfExpression(ifExp *ast.IfExpression) object.Object {
	cond := Eval(ifExp.Condition)

	if object.IsTruthy(cond) {
		return Eval(ifExp.Consequence)
	}
	if ifExp.Alternative != nil {
		return Eval(ifExp.Alternative)
	}
	return object.NULL
}

func evalBlockStmt(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt)

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
