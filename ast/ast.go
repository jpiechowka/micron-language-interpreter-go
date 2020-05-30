package ast

type AstNode interface {
	TokenLiteral() string
}

type Statement interface {
	AstNode
	statementNode()
}

type Expression interface {
	AstNode
	expressionNode()
}
