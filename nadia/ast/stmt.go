package ast

import (
	"github.com/agustin-del-pino/nadia-lang/nadia/ast/node"
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer/token"
)

type (
	Block struct {
		Op    *token.Tok
		Cl    *token.Tok
		Stmts []Stmt
	}
	StmtWrapper struct {
		Node Node
	}
	Return struct {
		Tk *token.Tok
		X  Expr
	}
	Assign struct {
		Op *token.Tok
		Lt Expr
		Rt Expr
	}
	If struct {
		Tk   *token.Tok
		Cond Expr
		Body *Block
		Else *If
	}
	For struct {
		Tk   *token.Tok
		Cond Expr
		Body *Block
	}
	ForRange struct {
		Tk      *token.Tok
		Counter *Ident
		Start   Expr
		Stop    Expr
		Step    Expr
		Body    *Block
	}
	ForEach struct {
		Tk       *token.Tok
		Holder   *Ident
		Iterator Expr
		Body     *Block
	}
	When struct {
		Tk        *token.Tok
		Cl        *token.Tok
		Val       Expr
		Clauses   []Expr
		Bodies    []*Block
		Otherwise *Block
	}
	FlowControl struct {
		Tk *token.Tok
	}
)

func (n *Block) Kind() node.Kind { return node.Block }
func (n *Block) St() int         { return n.Op.St }
func (n *Block) Ed() int         { return n.Cl.Ed }

func (n *StmtWrapper) Kind() node.Kind { return node.StmtWrapper }
func (n *StmtWrapper) St() int         { return n.Node.St() }
func (n *StmtWrapper) Ed() int         { return n.Node.Ed() }

func (n *Return) Kind() node.Kind { return node.Return }
func (n *Return) St() int         { return n.Tk.St }
func (n *Return) Ed() int {
	if n.X == nil {
		return n.Tk.Ed
	}
	return n.X.Ed()
}

func (n *Assign) Kind() node.Kind { return node.Assign }
func (n *Assign) St() int         { return n.Lt.St() }
func (n *Assign) Ed() int         { return n.Rt.Ed() }

func (n *If) Kind() node.Kind { return node.If }
func (n *If) St() int         { return n.Tk.St }
func (n *If) Ed() int {
	if n.Else != nil {
		return n.Else.Ed()
	}
	return n.Body.Ed()
}

func (n *For) Kind() node.Kind { return node.For }
func (n *For) St() int         { return n.Tk.St }
func (n *For) Ed() int         { return n.Body.Ed() }

func (n *ForRange) Kind() node.Kind { return node.ForRange }
func (n *ForRange) St() int         { return n.Tk.St }
func (n *ForRange) Ed() int         { return n.Body.Ed() }

func (n *ForEach) Kind() node.Kind { return node.ForEach }
func (n *ForEach) St() int         { return n.Tk.St }
func (n *ForEach) Ed() int         { return n.Body.Ed() }

func (n *When) Kind() node.Kind { return node.When }
func (n *When) St() int         { return n.Tk.St }
func (n *When) Ed() int         { return n.Cl.Ed }

func (n *FlowControl) Kind() node.Kind { return node.FlowControl }
func (n *FlowControl) St() int         { return n.Tk.St }
func (n *FlowControl) Ed() int         { return n.Tk.Ed }

func (n *Block) stmt()       {}
func (n *StmtWrapper) stmt() {}
func (n *Return) stmt()      {}
func (n *Assign) stmt()      {}
func (n *If) stmt()          {}
func (n *For) stmt()         {}
func (n *ForRange) stmt()    {}
func (n *ForEach) stmt()     {}
func (n *When) stmt()        {}
func (n *FlowControl) stmt() {}
