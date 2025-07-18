package eval

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		if !testIntegerObject(t, testEval(tt.input), tt.expected) {
			return
		}
	}
}

func testEval(input string) object.Object {
	return Eval(parser.New(lexer.New(input)).ParseProgram())
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not an Integer, got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object was expected to have value of %d, got=%d",
			expected, result.Value)
		return false
	}
	return true
}
