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
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2*2*2", 32},
		{"-50 +100 -50", 0},
		{"5* 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2* (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
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
	env := object.NewEnvironment()

	return Eval(program, env)
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
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
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

func TestBangOpp(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		eval := testEval(tt.input)
		testBoolObj(t, eval, tt.expected)
	}
}

func TestIfElse(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		eval := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntObj(t, eval, int64(integer))
		} else {
			testNullObj(t, eval)
		}
	}
}

func testNullObj(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("fail")
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9", 10},
		{"9; return 10; 9;", 10},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				}

				return 1;
			}
			`,
			10,
		},
	}

	for _, tt := range tests {
		eval := testEval(tt.input)
		testIntObj(t, eval, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}

				return 1;
			}
			`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - " " - "World"`,
			"unknown operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		eval := testEval(tt.input)

		errObj, ok := eval.(*object.Error)
		if !ok {
			t.Errorf("fail")
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("fail")
		}
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input string
		exp   int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5*5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntObj(t, testEval(tt.input), tt.exp)
	}
}

func TestFunc(t *testing.T) {
	input := "fn(x) { x + 2; };"
	eval := testEval(input)
	fn, ok := eval.(*object.Function)
	if !ok {
		t.Fatalf("fail")
	}

	if len(fn.Params) != 1 {
		t.Fatalf("fail")
	}

	if fn.Params[0].String() != "x" {
		t.Fatalf("Fail")
	}

	expBody := "(x + 2)"

	if fn.Body.String() != expBody {
		t.Fatalf("fail")
	}
}

func TestFuncApplication(t *testing.T) {
	tests := []struct {
		input string
		exp   int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5);", 5},
	}

	for _, tt := range tests {
		testIntObj(t, testEval(tt.input), tt.exp)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	eval := testEval(input)
	str, ok := eval.(*object.String)
	if !ok {
		t.Fatalf("fail")
	}
	if str.Value != "Hello World!" {
		t.Errorf("fail")
	}
}

func TestStringConcat(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	eval := testEval(input)
	str, ok := eval.(*object.String)
	if !ok {
		t.Fatalf("fail")
	}
	if str.Value != "Hello World!" {
		t.Errorf("fail")
	}
}

func TestBuiltin(t *testing.T) {
	tests := []struct {
		input string
		exp   interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}

	for _, tt := range tests {
		eval := testEval(tt.input)

		switch exp := tt.exp.(type) {
		case int:
			testIntObj(t, eval, int64(exp))
		case string:
			errObj, ok := eval.(*object.Error)
			if !ok {
				t.Errorf("fail")
				continue
			}
			if errObj.Message != exp {
				t.Errorf("fail")
			}
		}
	}
}

func TestArrayLit(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	eval := testEval(input)
	result, ok := eval.(*object.Array)
	if !ok {
		t.Fatalf("fail")
	}

	if len(result.Elements) != 3 {
		t.Fatalf("fail")
	}

	testIntObj(t, result.Elements[0], 1)
	testIntObj(t, result.Elements[1], 4)
	testIntObj(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntObj(t, evaluated, int64(integer))
		} else {
			testNullObj(t, evaluated)
		}
	}
}
