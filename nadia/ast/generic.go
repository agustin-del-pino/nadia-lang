package ast

import (
	"github.com/agustin-del-pino/nadia-lang/nadia/ast/node"
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer/token"
)

type (
	Source struct {
		Bg         *token.Tok
		Fs         *token.Tok
		Includes   []*Include
		Types      []*TypeDef
		Funcs      []*Func
		Vals       []*ValDef
		Defs       []*Def
		Objs       []*Obj
		Events     []*Event
		EventsFunc map[string][]*EventFunc
	}

	NodeList[T Node] struct {
		Op   *token.Tok
		Cl   *token.Tok
		List []T
	}

	TypedField struct {
		Name    *Ident
		Type    *Ident
		Pointer bool
	}
	ValueField struct {
		Name *Ident
		Val  Expr
	}
	OperationField struct {
		Op       *token.Tok
		Cl       *token.Tok
		Nodes    *LitVal
		Operator *token.Tok
	}
)

func (n *Source) Kind() node.Kind { return node.Source }
func (n *Source) St() int         { return n.Bg.St }
func (n *Source) Ed() int         { return n.Fs.Ed }

func (n *TypedField) Kind() node.Kind { return node.TypedField }
func (n *TypedField) St() int         { return n.Name.St() }
func (n *TypedField) Ed() int         { return n.Type.Ed() }

func (n *ValueField) Kind() node.Kind { return node.ValueField }
func (n *ValueField) St() int         { return n.Name.St() }
func (n *ValueField) Ed() int         { return n.Val.Ed() }

func (n *OperationField) Kind() node.Kind { return node.OperationField }
func (n *OperationField) St() int         { return n.Op.St }
func (n *OperationField) Ed() int         { return n.Cl.Ed }
