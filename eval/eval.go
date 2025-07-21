package eval

import (
	"monkey/ast"
	"monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		return evalPrefixExp(node)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return toBoolean(node.Value)
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

func toBoolean(val bool) *object.Boolean {
	if val {
		return TRUE
	}
	return FALSE
}

func evalPrefixExp(node *ast.PrefixExpression) object.Object {
	val := Eval(node.Right)

	switch node.Operator {
	case "!":
		return evalBang(val)
	case "-":
		return evalMinus(val)
	default:
		return NULL
	}
}

func evalBang(right object.Object) object.Object {
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

func evalMinus(right object.Object) object.Object {
	if right.Type() != object.OBJ_INTEGER {
		return NULL
	}
	integer := right.(*object.Integer)
	return &object.Integer{Value: -integer.Value}
}
