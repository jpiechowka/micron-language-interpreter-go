package parser

import (
	"fmt"
	"github.com/jpiechowka/micron-language-interpreter-go/ast"
	"github.com/jpiechowka/micron-language-interpreter-go/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let test_test = 13371337;
`

	parser := New(lexer.New(input))

	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if program == nil {
		t.Fatalf("ParseProgram() retruned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements length is not equal to 3. Actual length is: %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"test_test"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 13371337;
`

	parser := New(lexer.New(input))

	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements length is not equal to 3. Actual length is: %d", len(program.Statements))
	}

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("returnStatement is not *ast.ReturnStatement. Got %T instead", statement)
			continue
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStatement.TokenLiteral is not 'return'. Got %s instead", returnStatement.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	parser := New(lexer.New(input))

	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program hasn't got enough statements. Got %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. Got %T instead", program.Statements[0])
	}

	identifier, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("statement.Expression is not *ast.Identifier. Got %T instead", statement.Expression)
	}

	if identifier.Value != "foobar" {
		t.Errorf("identifier.Value is not %s. Got %s instead", "foobar", identifier.Value)
	}

	if identifier.TokenLiteral() != "foobar" {
		t.Errorf("identifier.TokenLiteral() is not %s. Got %s instead", "foobar", identifier.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	parser := New(lexer.New(input))

	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program hasn't got enough statements. Got %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. Got %T instead", program.Statements[0])
	}

	literal, ok := statement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("statement.Expression is not *ast.IntegerLiteral. Got %T instead", statement.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value is not %d. Got %d instead", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral is not %s. Got %s instead", "5", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!13371337;", "!", 13371337},
		{"-13371337;", "-", 13371337},
	}

	for _, prefixTest := range prefixTests {
		parser := New(lexer.New(prefixTest.input))

		program := parser.ParseProgram()
		checkParseErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. Got %d instead", 1, len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. Got %T instead", program.Statements[0])
		}

		expression, ok := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("statement.Expression is not ast.PrefixExpression. Got %T instead", statement.Expression)
		}

		if expression.Operator != prefixTest.operator {
			t.Fatalf("expression.Operator is not %s. Got %s instead", prefixTest.operator, expression.Operator)
		}

		if !testIntegerLiteral(t, expression.Right, prefixTest.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"77 + 1337;", 77, "+", 1337},
		{"77 - 1337;", 77, "-", 1337},
		{"77 * 1337;", 77, "*", 1337},
		{"77 / 1337;", 77, "/", 1337},
		{"77 > 1337;", 77, ">", 1337},
		{"77 < 1337;", 77, "<", 1337},
		{"77 == 1337;", 77, "==", 1337},
		{"77 != 1337;", 77, "!=", 1337},
	}

	for _, infixTest := range infixTests {
		parser := New(lexer.New(infixTest.input))

		program := parser.ParseProgram()
		checkParseErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. Got %d instead", 1, len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. Got %T instead", program.Statements[0])
		}

		expression, ok := statement.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("statement.Expression is not ast.InfixExpression. Got %T instead", statement.Expression)
		}

		if !testIntegerLiteral(t, expression.Left, infixTest.leftValue) {
			return
		}

		if expression.Operator != infixTest.operator {
			t.Fatalf("expression.Operator is not %s. Got %s instead", infixTest.operator, expression.Operator)
		}

		if !testIntegerLiteral(t, expression.Right, infixTest.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, precedenceTest := range tests {
		parser := New(lexer.New(precedenceTest.input))

		program := parser.ParseProgram()
		checkParseErrors(t, parser)

		actual := program.String()

		if actual != precedenceTest.expected {
			t.Errorf("Precedence test error - actual value %s is not the same as expected value: %s", actual, precedenceTest.expected)
		}
	}
}

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("statement.TokenLiteral is not let. Got %s instead", statement.TokenLiteral())
		return false
	}

	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("letStatement is not *ast.LetStatement. Got %T instead", statement)
		return false
	}

	if letStatement.Name.Value != name {
		t.Errorf("letStatement.Name.Value is not %s. Got %s instead", name, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("letStatement.Name.TokenLiteral() is not %s. Got %s instead", name, letStatement.Name.TokenLiteral())
		return false
	}

	return true
}

func checkParseErrors(t *testing.T, parser *Parser) {
	parserErrors := parser.GetErrors()
	if len(parserErrors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(parserErrors))
	for _, errorMsg := range parserErrors {
		t.Errorf("Parser error: %s", errorMsg)
	}

	t.FailNow()
}

func testIntegerLiteral(t *testing.T, integerLiteral ast.Expression, value int64) bool {
	integ, ok := integerLiteral.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("integerLiteral is not *ast.IntegerLiteral. Got %T instead", integerLiteral)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value is not %d. Got %d instead", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral() is not %d. Got %s instead", value,
			integ.TokenLiteral())
		return false
	}
	return true
}

func testIdentifier(t *testing.T, expression ast.Expression, value string) bool {
	identifier, ok := expression.(*ast.Identifier)
	if !ok {
		t.Errorf("expression is not *ast.Identifier. Got %T instead", expression)
		return false
	}

	if identifier.Value != value {
		t.Errorf("identifier.Value is not %s. Got %s instead", value, identifier.Value)
		return false
	}

	if identifier.TokenLiteral() != value {
		t.Errorf("identifier.TokenLiteral() is not %s. Got %s instead", value, identifier.TokenLiteral())
		return false
	}

	return true
}
