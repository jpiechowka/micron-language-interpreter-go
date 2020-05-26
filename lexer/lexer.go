package lexer

import "github.com/jpiechowka/micron-language-interpreter-go/token"

type Lexer struct {
	input            string // input string containing code
	currentChar      byte   // current char under examination
	currentPosition  int    // current position in input (points to the current char)
	nextReadPosition int    // current reading position in input (after current char)
}

func New(lexerInput string) *Lexer {
	lexer := &Lexer{input: lexerInput}
	lexer.readChar()
	return lexer
}

func (lexer *Lexer) NextToken() token.Token {
	var nextToken token.Token

	lexer.consumeWhitespaces()

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
	default:
		if isLetter(lexer.currentChar) {
			nextToken.Literal = lexer.readIdentifier()
			nextToken.TokenType = token.LookupIdentifier(nextToken.Literal)
			return nextToken
		} else if isDigit(lexer.currentChar) {
			nextToken.Literal = lexer.readInteger()
			nextToken.TokenType = token.INT
			return nextToken
		} else {
			nextToken = newToken(token.ILLEGAL, lexer.currentChar)
		}
	}

	lexer.readChar()
	return nextToken
}

func (lexer *Lexer) readChar() {
	if lexer.nextReadPosition >= len(lexer.input) { // Check if end of input is reached
		lexer.currentChar = 0 // 0 is NULL in ASCII
	} else {
		lexer.currentChar = lexer.input[lexer.nextReadPosition]
	}

	lexer.currentPosition = lexer.nextReadPosition
	lexer.nextReadPosition += 1
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{TokenType: tokenType, Literal: string(char)}
}

func (lexer *Lexer) readIdentifier() string {
	startPosition := lexer.currentPosition

	for isLetter(lexer.currentChar) {
		lexer.readChar()
	}

	return lexer.input[startPosition:lexer.currentPosition]
}

func (lexer *Lexer) readInteger() string {
	startPosition := lexer.currentPosition

	for isDigit(lexer.currentChar) {
		lexer.readChar()
	}

	return lexer.input[startPosition:lexer.currentPosition]
}

func isLetter(char byte) bool {
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || char == '_'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func (lexer *Lexer) consumeWhitespaces() {
	for lexer.currentChar == ' ' || lexer.currentChar == '\t' || lexer.currentChar == '\n' || lexer.currentChar == '\r' {
		lexer.readChar()
	}
}
