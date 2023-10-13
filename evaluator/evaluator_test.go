package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvalIntExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		eval := testEval(tt.input)
		testIntObj(t, eval, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntObj(t *testing.T, obj object.Object, exp int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer got %T (%+v)", obj, obj)
		return false
	}

	if result.Value != exp {
		t.Errorf("Object has wrong value. got %d want %d", result.Value, exp)
		return false
	}
	return true
}

func TestEvalBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		eval := testEval(tt.input)
		testBoolObj(t, eval, tt.expected)
	}
}

func testBoolObj(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("fail")
		return false
	}
	if result.Value != expected {
		t.Errorf("fail")
		return false
	}
	return true
}
