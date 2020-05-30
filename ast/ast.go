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

type Program struct {
	Statements []Statement
}

func (program *Program) TokenLiteral() string {
	if len(program.Statements) > 0 {
		return program.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
