package lexer

import (
	"github.com/jpiechowka/micron-language-interpreter-go/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
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
