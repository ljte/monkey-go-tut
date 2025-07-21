package eval

import (
	"monkey/ast"
	"monkey/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		return evalPrefixExp(node)
	case *ast.InfixExpression:
		return evalInfixExp(node)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return object.AsBool(node.Value)
	default:
		return nil
	}
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt)
	}
	return result
}

func evalPrefixExp(node *ast.PrefixExpression) object.Object {
	val := Eval(node.Right)

	switch node.Operator {
	case "!":
		return evalBang(val)
	case "-":
		return evalMinus(val)
	default:
		return object.NULL
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
		return object.NULL
	}
	integer := right.(*object.Integer)
	return object.AsInt(-integer.Value)
}

func evalInfixExp(node *ast.InfixExpression) object.Object {
	left := Eval(node.Left)
	right := Eval(node.Right)

	switch {
	case left.Type() == object.OBJ_INTEGER && right.Type() == object.OBJ_INTEGER:
		return evalIntegerInfixExp(left, right, node.Operator)
	case node.Operator == "==":
		return object.AsBool(left == right)
	case node.Operator == "!=":
		return object.AsBool(left != right)
	default:
		return object.NULL

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
		return object.NULL
	}
}
