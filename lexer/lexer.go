package lexer

import "github.com/jpiechowka/micron-language-interpreter-go/token"

type Lexer struct {
	lexerInput       string // input string containing code
	currentChar      byte   // current char under examination
	currentPosition  int    // current position in input (points to the current char)
	nextReadPosition int    // current reading position in input (after current char)
}

func New(lexerInput string) *Lexer {
	lexer := &Lexer{lexerInput: lexerInput}
	lexer.readChar()
	return lexer
}

func (lexer *Lexer) NextToken() token.Token {
	var nextToken token.Token

	switch lexer.currentChar {
	case '=':
		nextToken = newToken(token.ASSIGN, lexer.currentChar)
	case ';':
		nextToken = newToken(token.SEMICOLON, lexer.currentChar)
	case '(':
		nextToken = newToken(token.LPAREN, lexer.currentChar)
	case ')':
		nextToken = newToken(token.RPAREN, lexer.currentChar)
	case ',':
		nextToken = newToken(token.COMMA, lexer.currentChar)
	case '+':
		nextToken = newToken(token.PLUS, lexer.currentChar)
	case '{':
		nextToken = newToken(token.LBRACE, lexer.currentChar)
	case '}':
		nextToken = newToken(token.RBRACE, lexer.currentChar)
	case 0:
		nextToken.Literal = ""
		nextToken.TokenType = token.EOF
	}

	lexer.readChar()
	return nextToken
}

func (lexer *Lexer) readChar() {
	if lexer.nextReadPosition >= len(lexer.lexerInput) { // Check if end of input is reached
		lexer.currentChar = 0 // 0 is NULL in ASCII
	} else {
		lexer.currentChar = lexer.lexerInput[lexer.nextReadPosition]
	}

	lexer.currentPosition = lexer.nextReadPosition
	lexer.nextReadPosition += 1
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{TokenType: tokenType, Literal: string(char)}
}
