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

		let result = add(one, seven);
        !-/*10;
        0 < 10 >   5;

        if (     5 <   10) {
            return true;
        } else {
            return false;
        }

		7777 == 5555;
        1111 != 33333333;
`

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
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "0"},
		{token.LESSTHAN, "<"},
		{token.INT, "10"},
		{token.GREATERTHAN, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LESSTHAN, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "7777"},
		{token.EQUALITY, "=="},
		{token.INT, "5555"},
		{token.SEMICOLON, ";"},
		{token.INT, "1111"},
		{token.INEQUALITY, "!="},
		{token.INT, "33333333"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for i, tokenType := range tests {
		testedToken := lexer.NextToken()

		if testedToken.TokenType != tokenType.expectedType {
			t.Fatalf("Lexer test case [%d/%d] failed - TokenType is wrong. Expected %q, got %q", i, len(tests), tokenType.expectedType, testedToken.TokenType)
		}

		if testedToken.Literal != tokenType.expectedLiteral {
			t.Fatalf("Lexer tests case [%d/%d] failed - Literal is wrong. Expected %q, got %q", i, len(tests), tokenType.expectedLiteral, testedToken.Literal)
		}
	}
}
