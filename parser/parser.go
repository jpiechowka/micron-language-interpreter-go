package parser

import (
	"fmt"
	"github.com/jpiechowka/micron-language-interpreter-go/ast"
	"github.com/jpiechowka/micron-language-interpreter-go/lexer"
	"github.com/jpiechowka/micron-language-interpreter-go/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS        // ==
	LESSORGREATER // > or <
	SUM           // +
	PRODUCT       // *
	PREFIX        // -X or !X
	CALL          // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQUALITY:    EQUALS,
	token.INEQUALITY:  EQUALS,
	token.LESSTHAN:    LESSORGREATER,
	token.GREATERTHAN: LESSORGREATER,
	token.PLUS:        SUM,
	token.MINUS:       SUM,
	token.SLASH:       PRODUCT,
	token.ASTERISK:    PRODUCT,
}

type (
	prefixParseFunction func() ast.Expression
	infixParseFunction  func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer        *lexer.Lexer
	errors       []string
	currentToken token.Token
	peekToken    token.Token

	prefixParseFunctions map[token.TokenType]prefixParseFunction
	infixParseFunctions  map[token.TokenType]infixParseFunction
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	parser.prefixParseFunctions = make(map[token.TokenType]prefixParseFunction)

	parser.registerPrefix(token.IDENT, parser.parseIdentifier)
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefix(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)

	parser.infixParseFunctions = make(map[token.TokenType]infixParseFunction)

	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)
	parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfix(token.EQUALITY, parser.parseInfixExpression)
	parser.registerInfix(token.INEQUALITY, parser.parseInfixExpression)
	parser.registerInfix(token.LESSTHAN, parser.parseInfixExpression)
	parser.registerInfix(token.GREATERTHAN, parser.parseInfixExpression)

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
		return parser.parseExpressionStatement()
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

func (parser *Parser) registerPrefix(tokenType token.TokenType, function prefixParseFunction) {
	parser.prefixParseFunctions[tokenType] = function
}

func (parser *Parser) registerInfix(tokenType token.TokenType, function infixParseFunction) {
	parser.infixParseFunctions[tokenType] = function
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: parser.currentToken}

	statement.Expression = parser.parseExpression(LOWEST)

	if parser.isComparedTokenSameAsPeek(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefix := parser.prefixParseFunctions[parser.currentToken.TokenType]

	if prefix == nil {
		parser.noPrefixParseFunctionError(parser.currentToken.TokenType)
		return nil
	}

	leftExpression := prefix()

	for !parser.isComparedTokenSameAsPeek(token.SEMICOLON) && precedence < parser.peekPrecedence() {
		infix := parser.infixParseFunctions[parser.peekToken.TokenType]

		if infix == nil {
			return leftExpression
		}

		parser.nextToken()

		leftExpression = infix(leftExpression)
	}

	return leftExpression
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: parser.currentToken}

	value, err := strconv.ParseInt(parser.currentToken.Literal, 0, 64)

	if err != nil {
		errorMsg := fmt.Sprintf("Could not parse %q as integer", parser.currentToken.Literal)
		parser.errors = append(parser.errors, errorMsg)
		return nil
	}

	literal.Value = value

	return literal
}

func (parser *Parser) noPrefixParseFunctionError(tokenType token.TokenType) {
	errorMsg := fmt.Sprintf("No prefix parse function for %s found", tokenType)
	parser.errors = append(parser.errors, errorMsg)
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Literal,
	}

	parser.nextToken()

	expression.Right = parser.parseExpression(PREFIX)

	return expression
}

func (parser *Parser) peekPrecedence() int {
	if precedence, ok := precedences[parser.peekToken.TokenType]; ok {
		return precedence
	}

	return LOWEST
}

func (parser *Parser) currentPrecedence() int {
	if precedence, ok := precedences[parser.currentToken.TokenType]; ok {
		return precedence
	}

	return LOWEST
}

func (parser *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Literal,
		Left:     left,
	}

	precedence := parser.currentPrecedence()

	parser.nextToken()

	expression.Right = parser.parseExpression(precedence)

	return expression
}
