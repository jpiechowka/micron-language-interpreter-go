package lexer

import (
	"github.com/jpiechowka/micron-language-interpreter-go/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
		let one = 1;
		let seven = 7;

		let add = fn(x, y) {
			x + y;
		};

		let result = add(one, seven);`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "one"},
		{token.ASSIGN, "="},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "seven"},
		{token.ASSIGN, "="},
		{token.INT, "7"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "one"},
		{token.COMMA, ","},
		{token.IDENT, "seven"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for i, tokenType := range tests {
		testedToken := lexer.NextToken()

		if testedToken.TokenType != tokenType.expectedType {
			t.Fatalf("lexer test case [%d/%d] TokenType is wrong. Expected %q, got %q", i, len(tests), tokenType.expectedType, testedToken.TokenType)
		}

		if testedToken.Literal != tokenType.expectedLiteral {
			t.Fatalf("lexer tests case [%d/%d] Literal is wrong. Expected %q, got %q", i, len(tests), tokenType.expectedLiteral, testedToken.Literal)
		}
	}
}
