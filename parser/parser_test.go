package parser

import (
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
