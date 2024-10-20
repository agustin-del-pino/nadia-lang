package parser

import (
	"github.com/agustin-del-pino/nadia-lang/nadia/ast"
	"github.com/agustin-del-pino/nadia-lang/nadia/ast/node"
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer"
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer/token"
)

var opAssign = set{}.Add(token.Assign).Add(token.AddAssign).Add(token.SubAssign).Add(token.MulAssign).Add(token.DivAssign).Add(token.PowAssign)

func parseReturn(t *lexer.Tokenizer, c *cursor) *ast.Return {
	c.assert(token.Return.String(), token.Return)
	n := ast.Return{
		Tk: c.token().Clone(),
	}
	c.next(t)

	if c.Kind == token.RCurvy {
		return &n
	}

	n.X = parseExpr(t, c)

	return &n
}

func parseWhen(t *lexer.Tokenizer, c *cursor) *ast.When {
	c.assert(token.When.String(), token.When)
	n := ast.When{
		Tk: c.token().Clone(),
	}
	c.next(t)

	n.Val = parseExpr(t, c)

	c.assert(token.LCurvy.String(), token.LCurvy)
	c.next(t)

	for c.Kind != token.RCurvy {
		c.assert(token.Is.String()+" or "+token.Otherwise.String(), token.Is, token.Otherwise)

		if c.Kind == token.Otherwise {
			c.next(t)
			n.Otherwise = parseBlock(t, c)
		} else {
			c.next(t)
			n.Clauses = append(n.Clauses, parseExpr(t, c))
			n.Bodies = append(n.Bodies, parseBlock(t, c))
		}
	}

	n.Cl = c.token().Clone()
	c.next(t)

	return &n
}

func parseFor(t *lexer.Tokenizer, c *cursor) ast.Stmt {
	c.assert(token.For.String(), token.For)
	tk := c.token().Clone()
	c.next(t)

	if c.Kind == token.LCurvy {
		return &ast.For{
			Tk:   tk,
			Body: parseBlock(t, c),
		}
	}

	x := parseExpr(t, c)

	if x.Kind() == node.Ident {
		if c.Kind == token.In {
			c.next(t)
			return &ast.ForEach{
				Tk:       tk,
				Holder:   x.(*ast.Ident),
				Iterator: parseExpr(t, c),
				Body:     parseBlock(t, c),
			}
		} else if c.Kind == token.Range {
			c.next(t)

			n := ast.ForRange{
				Tk:      tk,
				Counter: x.(*ast.Ident),
			}

			n.Start = parseExpr(t, c)
			c.assert(token.Comma.String(), token.Comma)
			c.next(t)
			n.Stop = parseExpr(t, c)

			if c.Kind == token.Comma {
				c.next(t)
				n.Step = parseExpr(t, c)
			}

			n.Body = parseBlock(t, c)

			return &n
		}
	}

	return &ast.For{
		Tk:   tk,
		Cond: x,
		Body: parseBlock(t, c),
	}
}

func parseIf(t *lexer.Tokenizer, c *cursor) *ast.If {
	c.assert(token.If.String(), token.If)
	n := ast.If{
		Tk: c.token().Clone(),
	}
	c.next(t)

	n.Cond = parseExpr(t, c)
	n.Body = parseBlock(t, c)

	if c.Kind == token.Else {
		tk := c.token().Clone()
		c.next(t)
		if c.Kind == token.If {
			n.Else = parseIf(t, c)
			n.Else.Tk = tk
		} else {
			n.Else = &ast.If{
				Tk:   tk,
				Body: parseBlock(t, c),
			}
		}
	}
	return &n
}

func parseBlock(t *lexer.Tokenizer, c *cursor) *ast.Block {
	c.assert(token.LCurvy.String(), token.LCurvy)
	n := ast.Block{
		Op: c.token().Clone(),
	}
	c.next(t)

	if c.Kind == token.RCurvy {
		n.Cl = c.token().Clone()
		c.next(t)
		return &n
	}

	for c.Kind != token.RCurvy {
		n.Stmts = append(n.Stmts, parseStmt(t, c))
	}

	n.Cl = c.token().Clone()
	c.next(t)

	return &n
}

func parseFlowControl(t *lexer.Tokenizer, c *cursor) *ast.FlowControl {
	c.assert(token.Stop.String()+"or"+token.Pass.String(), token.Stop, token.Pass)
	defer c.next(t)
	return &ast.FlowControl{
		Tk: c.token().Clone(),
	}
}

func parseStmt(t *lexer.Tokenizer, c *cursor) ast.Stmt {
	switch c.Kind {
	case token.Var, token.Const:
		return &ast.StmtWrapper{Node: parseValDef(t, c)}
	case token.Return:
		return parseReturn(t, c)
	case token.If:
		return parseIf(t, c)
	case token.For:
		return parseFor(t, c)
	case token.When:
		return parseWhen(t, c)
	case token.Stop, token.Pass:
		return parseFlowControl(t, c)
	default:
		x := parseExpr(t, c)

		if x.Kind() == node.Call {
			return &ast.StmtWrapper{Node: x}
		} else if opAssign.Has(c.Kind) {
			op := c.token().Clone()
			c.next(t)
			return &ast.Assign{
				Lt: x,
				Op: op,
				Rt: parseExpr(t, c),
			}
		}
		panic("invalid statement: " + c.Val)
	}
}
