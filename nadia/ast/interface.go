package ast

import "github.com/agustin-del-pino/nadia-lang/nadia/ast/node"

type Node interface {
	Kind() node.Kind
	St() int
	Ed() int
}

type Expr interface {
	Node
	expr()
}

type Stmt interface {
	Node
	stmt()
}

type Decl interface {
	Node
	decl()
}
