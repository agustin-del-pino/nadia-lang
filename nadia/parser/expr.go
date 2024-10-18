package parser

import (
	"github.com/agustin-del-pino/nadia-lang/nadia/ast"
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer"
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer/token"
)

var opBinaryOr = set{}.Add(token.Or)
var opBinaryAnd = set{}.Add(token.And)
var opBinaryEquality = set{}.Add(token.Equal).Add(token.NotEqual)
var opBinaryRelational = set{}.Add(token.Less).Add(token.LEqual).Add(token.Greater).Add(token.GEqual)
var opBinaryTerm = set{}.Add(token.Add).Add(token.Sub)
var opBinaryFactor = set{}.Add(token.Mul).Add(token.Div).Add(token.Mod)
var opUnary = set{}.Add(token.Sub).Add(token.Not)

func parseIdent(t *lexer.Tokenizer, c *cursor) *ast.Ident {
	c.assert("identifier", token.Ident)
	defer c.next(t)
	return &ast.Ident{
		Tk: c.token().Clone(),
	}
}

func parseLitVal(t *lexer.Tokenizer, c *cursor) *ast.LitVal {
	c.assert("literal value", token.String, token.Int, token.Binary, token.Hex, token.Float, token.Char)
	defer c.next(t)
	return &ast.LitVal{
		Tk: c.token().Clone(),
	}
}

func parseParenthesis(t *lexer.Tokenizer, c *cursor) *ast.Parenthesis {
	c.assert(token.LParen.String(), token.LParen)
	n := ast.Parenthesis{
		Op: c.token().Clone(),
	}
	c.next(t)

	n.X = parseExpr(t, c)

	c.assert(token.RParen.String(), token.RParen)
	n.Cl = c.token().Clone()
	c.next(t)

	return &n
}

func parseCall(t *lexer.Tokenizer, c *cursor, x ast.Expr) *ast.Call {
	return &ast.Call{
		Caller: x,
		Args:   parseNodeListExpr(t, c, token.LParen, token.RParen, token.Comma),
	}
}

func parseSelector(t *lexer.Tokenizer, c *cursor, x ast.Expr) *ast.Selector {
	c.assert(token.Dot.String(), token.Dot)
	c.next(t)
	return &ast.Selector{
		Path:   x,
		Select: parseIdent(t, c),
	}
}

func parseNew(t *lexer.Tokenizer, c *cursor) *ast.New {
	c.assert(token.New.String(), token.New)
	c.next(t)
	return &ast.New{
		Name:   parseIdent(t, c),
		Fields: parseNodeListValueField(t, c, token.LCurvy, token.RCurvy, token.Comma),
	}
}

func parseUnitExpr(t *lexer.Tokenizer, c *cursor) ast.Expr {
	var x ast.Expr
	switch c.Kind {
	case token.Ident:
		x = parseIdent(t, c)
	case token.LParen:
		x = parseParenthesis(t, c)
	case token.New:
		x = parseNew(t, c)
	default:
		x = parseLitVal(t, c)
	}
	for {
		if c.Kind == token.Dot {
			x = parseSelector(t, c, x)
		} else if c.Kind == token.LParen {
			x = parseCall(t, c, x)
		} else {
			break
		}
	}
	return x
}

func parseUnary(t *lexer.Tokenizer, c *cursor) ast.Expr {
	if !opUnary.Has(c.Kind) {
		return parseUnitExpr(t, c)
	}
	n := ast.Unary{
		Op: c.token().Clone(),
	}
	c.next(t)

	n.X = parseUnitExpr(t, c)

	return &n
}

func parseBinaryExpr(k set, t *lexer.Tokenizer, c *cursor, fn func(*lexer.Tokenizer, *cursor) ast.Expr) ast.Expr {
	n := fn(t, c)
	for k.Has(c.Kind) {
		tk := c.token().Clone()
		c.next(t)
		n = &ast.Binary{
			Lt: n,
			Rt: parseBinaryAnd(t, c),
			Op: tk,
		}
	}
	return n
}

func parseBinaryFactor(t *lexer.Tokenizer, c *cursor) ast.Expr {
	return parseBinaryExpr(opBinaryFactor, t, c, parseUnary)
}

func parseBinaryTerm(t *lexer.Tokenizer, c *cursor) ast.Expr {
	return parseBinaryExpr(opBinaryTerm, t, c, parseBinaryFactor)
}

func parseBinaryRelational(t *lexer.Tokenizer, c *cursor) ast.Expr {
	return parseBinaryExpr(opBinaryRelational, t, c, parseBinaryTerm)
}

func parseBinaryEquality(t *lexer.Tokenizer, c *cursor) ast.Expr {
	return parseBinaryExpr(opBinaryEquality, t, c, parseBinaryRelational)
}

func parseBinaryAnd(t *lexer.Tokenizer, c *cursor) ast.Expr {
	return parseBinaryExpr(opBinaryAnd, t, c, parseBinaryEquality)
}

func parseBinaryOr(t *lexer.Tokenizer, c *cursor) ast.Expr {
	return parseBinaryExpr(opBinaryOr, t, c, parseBinaryAnd)
}

func parseExpr(t *lexer.Tokenizer, c *cursor) ast.Expr {
	return parseBinaryOr(t, c)
}
