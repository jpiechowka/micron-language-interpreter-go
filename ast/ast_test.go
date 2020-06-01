package ast

import (
	"github.com/jpiechowka/micron-language-interpreter-go/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{TokenType: token.LET, Literal: "let"},
				Name:  &Identifier{Token: token.Token{TokenType: token.IDENT, Literal: "someVar"}, Value: "someVar"},
				Value: &Identifier{Token: token.Token{TokenType: token.IDENT, Literal: "anotherVar"}, Value: "anotherVar"},
			},
		},
	}

	if program.String() != "let someVar = anotherVar;" {
		t.Errorf("AST test string failed - program.String() returned wrong values. Got %s", program.String())
	}
}
