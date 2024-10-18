package parser

import (
	"github.com/agustin-del-pino/nadia-lang/nadia/ast"
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer"
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer/token"
)

func parseTypeDef(t *lexer.Tokenizer, c *cursor) *ast.TypeDef {
	c.assert(token.Type.String(), token.Type)
	n := ast.TypeDef{
		Tk: c.token().Clone(),
	}
	c.next(t)
	n.Name = parseIdent(t, c)
	c.assert(token.Arrow.String(), token.Arrow)
	c.next(t)
	n.Default = c.token().Clone()
	c.next(t)
	c.assert(token.Arrow.String(), token.Arrow)
	c.next(t)
	n.Ref = parseLitVal(t, c)
	n.Ops = parseNodeListOperationField(t, c, token.LCurvy, token.RCurvy, token.Comma)
	return &n
}

func parseDef(t *lexer.Tokenizer, c *cursor) *ast.Def {
	c.assert(token.Def.String(), token.Def)
	n := ast.Def{
		Tk: c.token().Clone(),
	}
	c.next(t)

	n.Name = parseIdent(t, c)

	if c.Kind == token.LParen {
		n.Args = parseNodeListIdent(t, c, token.LParen, token.RParen, token.Comma)
	}

	if c.Kind == token.Ident {
		n.Type = parseIdent(t, c)
	}

	c.assert(token.Assign.String(), token.Assign)
	c.next(t)

	n.Val = parseLitVal(t, c)

	return &n
}

func parseValDef(t *lexer.Tokenizer, c *cursor) *ast.ValDef {
	c.assert(token.Var.String()+" or "+token.Const.String(), token.Var, token.Const)
	n := ast.ValDef{
		Tk: c.token().Clone(),
	}
	c.next(t)

	n.Name = parseIdent(t, c)

	if c.Kind == token.Ident {
		n.Type = parseIdent(t, c)
	}

	c.assert(token.Assign.String(), token.Assign)
	c.next(t)

	n.Val = parseExpr(t, c)

	return &n
}

func parseObj(t *lexer.Tokenizer, c *cursor) *ast.Obj {
	c.assert(token.Obj.String(), token.Obj)
	n := ast.Obj{
		Tk: c.token().Clone(),
	}
	c.next(t)

	n.Name = parseIdent(t, c)
	n.Props = parseNodeListTypedField(t, c, token.LCurvy, token.RCurvy, token.Comma)

	return &n
}

func parseEvent(t *lexer.Tokenizer, c *cursor) *ast.Event {
	c.assert(token.Event.String(), token.Event)
	n := ast.Event{
		Tk: c.token().Clone(),
	}
	c.next(t)

	n.Name = parseIdent(t, c)

	c.assert(token.Arrow.String(), token.Arrow)
	c.next(t)

	n.Parent = parseIdent(t, c)
	n.Props = parseNodeListTypedField(t, c, token.LCurvy, token.RCurvy, token.Comma)

	return &n
}

func parseFunc(t *lexer.Tokenizer, c *cursor) *ast.Func {
	c.assert(token.Func.String(), token.Func)

	n := ast.Func{
		Tk: c.token().Clone(),
	}
	c.next(t)

	n.Name = parseIdent(t, c)
	n.Params = parseNodeListTypedField(t, c, token.LParen, token.RParen, token.Comma)

	if c.Kind == token.Ident {
		n.Type = parseIdent(t, c)
	}

	n.Body = parseBlock(t, c)

	return &n
}

func parseEventFunc(t *lexer.Tokenizer, c *cursor) *ast.EventFunc {
	c.assert(token.Ev.String()+" or "+token.Lst.String(), token.Ev, token.Lst)

	n := ast.EventFunc{
		Tk: c.token().Clone(),
	}

	c.next(t)

	c.assert(token.LParen.String(), token.LParen)
	c.next(t)

	n.Event = parseIdent(t, c)

	c.assert(token.RParen.String(), token.RParen)
	c.next(t)

	if c.Kind == token.Func {
		n.Func = parseFunc(t, c)
	} else {
		n.Link = parseIdent(t, c)
	}

	return &n
}

func parseInclude(t *lexer.Tokenizer, c *cursor) *ast.Include {
	c.assert(token.Include.String(), token.Include)
	n := ast.Include{
		Tk: c.token().Clone(),
	}
	c.next(t)

	n.Path = parseLitVal(t, c)

	if c.Kind == token.As {
		c.next(t)
		n.Alias = parseIdent(t, c)
	}
	return &n
}
