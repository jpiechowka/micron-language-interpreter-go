package parser

import (
	"fmt"
	"github.com/jpiechowka/micron-language-interpreter-go/ast"
	"github.com/jpiechowka/micron-language-interpreter-go/lexer"
	"github.com/jpiechowka/micron-language-interpreter-go/token"
)

type (
	prefixParseFunction func() ast.Expression
	infixParseFunction  func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer        *lexer.Lexer
	errors       []string
	currentToken token.Token
	peekToken    token.Token
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:  lexer,
		errors: []string{},
	}
	parser.nextToken() // Read two tokens so currentToken and peekToken are populated
	parser.nextToken()
	return parser
}

func (parser *Parser) GetErrors() []string {
	return parser.errors
}

func (parser *Parser) nextToken() {
	parser.currentToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}

	for parser.currentToken.TokenType != token.EOF {
		statement := parser.parseStatement()

		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		parser.nextToken()
	}

	return program
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.currentToken.TokenType {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return nil
	}
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: parser.currentToken}

	if !parser.expectPeek(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}

	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: We're skipping the expressions until we encounter a semicolon
	for !parser.isComparedTokenSameAsCurrent(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: parser.currentToken}

	parser.nextToken()

	// TODO: We're skipping the expressions until we encounter a semicolon
	for !parser.isComparedTokenSameAsCurrent(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) isComparedTokenSameAsCurrent(tokenToCompare token.TokenType) bool {
	return parser.currentToken.TokenType == tokenToCompare
}

func (parser *Parser) isComparedTokenSameAsPeek(tokenToCompare token.TokenType) bool {
	return parser.peekToken.TokenType == tokenToCompare
}

func (parser *Parser) expectPeek(expectedToken token.TokenType) bool {
	if parser.isComparedTokenSameAsPeek(expectedToken) {
		parser.nextToken()
		return true
	} else {
		parser.peekError(expectedToken)
		return false
	}
}

func (parser *Parser) peekError(token token.TokenType) {
	errorMsg := fmt.Sprintf("Expected next token to be %s, got %s instead", token, parser.peekToken.TokenType)
	parser.errors = append(parser.errors, errorMsg)
}
