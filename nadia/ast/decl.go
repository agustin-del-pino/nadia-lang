package ast

import (
	"github.com/agustin-del-pino/nadia-lang/nadia/ast/node"
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer/token"
)

type (
	TypeDef struct {
		Tk      *token.Tok
		Name    *Ident
		Default *token.Tok
		Ref     *LitVal
		Ops     *NodeList[*OperationField]
	}
	Def struct {
		Tk   *token.Tok
		Name *Ident
		Type *Ident
		Args *NodeList[*Ident]
		Val  *LitVal
	}
	ValDef struct {
		Tk   *token.Tok
		Name *Ident
		Type *Ident
		Val  Expr
	}
	Obj struct {
		Tk    *token.Tok
		Name  *Ident
		Props *NodeList[*TypedField]
	}
	Func struct {
		Tk     *token.Tok
		Name   *Ident
		Params *NodeList[*TypedField]
		Type   *Ident
		Body   *Block
	}
	Event struct {
		Tk     *token.Tok
		Name   *Ident
		Parent *Ident
		Props  *NodeList[*TypedField]
	}
	EventFunc struct {
		Tk    *token.Tok
		Event *Ident
		Link  *Ident
		Func  *Func
	}
	Include struct {
		Tk    *token.Tok
		Path  *LitVal
		Alias *Ident
	}
)

func (n *TypeDef) Kind() node.Kind { return node.TypeDef }
func (n *TypeDef) St() int         { return n.Tk.St }
func (n *TypeDef) Ed() int         { return n.Ops.Cl.Ed }

func (n *Def) Kind() node.Kind { return node.Definition }
func (n *Def) St() int         { return n.Tk.St }
func (n *Def) Ed() int         { return n.Val.Ed() }

func (n *ValDef) Kind() node.Kind { return node.ValueDef }
func (n *ValDef) St() int         { return n.Tk.St }
func (n *ValDef) Ed() int         { return n.Val.Ed() }

func (n *Obj) Kind() node.Kind { return node.Obj }
func (n *Obj) St() int         { return n.Tk.St }
func (n *Obj) Ed() int         { return n.Props.Cl.Ed }

func (n *Func) Kind() node.Kind { return node.Func }
func (n *Func) St() int         { return n.Tk.St }
func (n *Func) Ed() int         { return n.Body.Ed() }

func (n *Event) Kind() node.Kind { return node.EventDef }
func (n *Event) St() int         { return n.Tk.St }
func (n *Event) Ed() int         { return n.Props.Cl.Ed }

func (n *EventFunc) Kind() node.Kind { return node.EventFunc }
func (n *EventFunc) St() int         { return n.Tk.St }
func (n *EventFunc) Ed() int {
	if n.Link != nil {
		return n.Link.Ed()
	}
	return n.Func.Ed()
}

func (n *Include) Kind() node.Kind { return node.Include }
func (n *Include) St() int         { return n.Tk.St }
func (n *Include) Ed() int {
	if n.Alias != nil {
		return n.Alias.Ed()
	}
	return n.Path.Ed()
}

func (n *TypeDef) decl()   {}
func (n *Def) decl()       {}
func (n *ValDef) decl()    {}
func (n *Obj) decl()       {}
func (n *Func) decl()      {}
func (n *Event) decl()     {}
func (n *EventFunc) decl() {}
func (n *Include) decl()   {}
