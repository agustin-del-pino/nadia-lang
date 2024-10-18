package parser

import (
	"github.com/agustin-del-pino/nadia-lang/nadia/ast"
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer"
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer/token"
)

func parseTypedField(t *lexer.Tokenizer, c *cursor) *ast.TypedField {
	n := ast.TypedField{
		Name: parseIdent(t, c),
	}
	if c.Kind == token.Mul {
		n.Pointer = true
		c.next(t)
	}
	n.Type = parseIdent(t, c)
	return &n
}
func parseValueField(t *lexer.Tokenizer, c *cursor) *ast.ValueField {
	n := ast.ValueField{
		Name: parseIdent(t, c),
	}
	c.assert(token.Colon.String(), token.Colon)
	c.next(t)
	n.Val = parseExpr(t, c)
	return &n
}

func parseOperationField(t *lexer.Tokenizer, c *cursor) *ast.OperationField {
	c.assert(token.LParen.String(), token.LParen)
	c.next(t)
	n := ast.OperationField{
		Nodes: parseLitVal(t, c),
	}
	c.assert(token.Comma.String(), token.Comma)
	c.next(t)
	n.Operator = c.token().Clone()
	c.next(t)
	c.assert(token.RParen.String(), token.RParen)
	c.next(t)
	return &n
}

func parseNodeList[T ast.Node](t *lexer.Tokenizer, c *cursor, op, cl, sp token.Kind, p func(t *lexer.Tokenizer, c *cursor) T) *ast.NodeList[T] {
	c.assert(op.String(), op)

	n := ast.NodeList[T]{
		Op: c.token().Clone(),
	}

	c.next(t)

	if c.Kind == cl {
		n.Cl = c.token().Clone()
		c.next(t)
		return &n
	}

	n.List = append(n.List, p(t, c))

	for c.Kind == sp {
		c.next(t)
		n.List = append(n.List, p(t, c))
	}

	c.assert(cl.String(), cl)
	n.Cl = c.token().Clone()
	c.next(t)
	return &n
}

func parseNodeListExpr(t *lexer.Tokenizer, c *cursor, op, cl, sp token.Kind) *ast.NodeList[ast.Expr] {
	return parseNodeList(t, c, op, cl, sp, parseExpr)
}

func parseNodeListIdent(t *lexer.Tokenizer, c *cursor, op, cl, sp token.Kind) *ast.NodeList[*ast.Ident] {
	return parseNodeList(t, c, op, cl, sp, parseIdent)
}

func parseNodeListTypedField(t *lexer.Tokenizer, c *cursor, op, cl, sp token.Kind) *ast.NodeList[*ast.TypedField] {
	return parseNodeList(t, c, op, cl, sp, parseTypedField)
}

func parseNodeListValueField(t *lexer.Tokenizer, c *cursor, op, cl, sp token.Kind) *ast.NodeList[*ast.ValueField] {
	return parseNodeList(t, c, op, cl, sp, parseValueField)
}

func parseNodeListOperationField(t *lexer.Tokenizer, c *cursor, op, cl, sp token.Kind) *ast.NodeList[*ast.OperationField] {
	return parseNodeList(t, c, op, cl, sp, parseOperationField)
}
