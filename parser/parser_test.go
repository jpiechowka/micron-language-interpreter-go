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
		t.Fatalf("Parser test failed - ParseProgram() retruned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Parser test failed - program.Statements length is not equal to 3. Actual length is: %d", len(program.Statements))
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

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("statement.TokenLiteral is not let. Got %s instead", statement.TokenLiteral())
		return false
	}

	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("statement is not *ast.LetStatement. Got %T instead", statement)
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
