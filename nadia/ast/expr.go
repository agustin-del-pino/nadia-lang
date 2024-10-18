package ast

import (
	"github.com/agustin-del-pino/nadia-lang/nadia/ast/node"
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer/token"
)

type (
	Ident struct {
		Tk      *token.Tok
	}
	LitVal struct {
		Tk *token.Tok
	}
	Unary struct {
		Op *token.Tok
		X  Expr
	}
	Binary struct {
		Op *token.Tok
		Lt Expr
		Rt Expr
	}
	Parenthesis struct {
		Op *token.Tok
		Cl *token.Tok
		X  Expr
	}
	Call struct {
		Caller Expr
		Args   *NodeList[Expr]
	}
	Selector struct {
		Path   Expr
		Select *Ident
	}
	New struct {
		Name   *Ident
		Fields *NodeList[*ValueField]
	}
)

func (n *Ident) Kind() node.Kind { return node.Ident }
func (n *Ident) St() int         { return n.Tk.St }
func (n *Ident) Ed() int         { return n.Tk.Ed }

func (n *LitVal) Kind() node.Kind { return node.LitVal }
func (n *LitVal) St() int         { return n.Tk.St }
func (n *LitVal) Ed() int         { return n.Tk.Ed }

func (n *Unary) Kind() node.Kind { return node.Unary }
func (n *Unary) St() int         { return n.Op.St }
func (n *Unary) Ed() int         { return n.X.Ed() }

func (n *Binary) Kind() node.Kind { return node.Binary }
func (n *Binary) St() int         { return n.Lt.St() }
func (n *Binary) Ed() int         { return n.Rt.Ed() }

func (n *Parenthesis) Kind() node.Kind { return node.Parenthesis }
func (n *Parenthesis) St() int         { return n.Op.St }
func (n *Parenthesis) Ed() int         { return n.Cl.St }

func (n *Call) Kind() node.Kind { return node.Call }
func (n *Call) St() int         { return n.Caller.St() }
func (n *Call) Ed() int         { return n.Args.Cl.Ed }

func (n *Selector) Kind() node.Kind { return node.Selector }
func (n *Selector) St() int         { return n.Path.St() }
func (n *Selector) Ed() int         { return n.Select.Ed() }

func (n *New) Kind() node.Kind { return node.New }
func (n *New) St() int         { return n.Name.St() }
func (n *New) Ed() int         { return n.Fields.Cl.Ed }

func (n *Ident) expr()       {}
func (n *LitVal) expr()      {}
func (n *Unary) expr()       {}
func (n *Binary) expr()      {}
func (n *Parenthesis) expr() {}
func (n *Call) expr()        {}
func (n *Selector) expr()    {}
func (n *New) expr()         {}
