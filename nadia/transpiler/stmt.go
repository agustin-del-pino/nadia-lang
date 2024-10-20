package transpiler

import (
	"bytes"

	"github.com/agustin-del-pino/nadia-lang/nadia/ast"
	"github.com/agustin-del-pino/nadia-lang/nadia/ast/node"
)

func transpileReturn(b *bytes.Buffer, scp *scope, n *ast.Return) {
	b.WriteString("return ")
	transpileExpr(b, scp, n.X)
	b.WriteByte(';')
}

func transpileStmtWrapper(b *bytes.Buffer, scp *scope, n *ast.StmtWrapper) {
	switch n.Node.Kind() {
	case node.ValueDef:
		transpileValDef(b, scp, n.Node.(*ast.ValDef))
	case node.Call:
		transpileCall(b, scp, n.Node.(*ast.Call))
		b.WriteByte(';')
	}
}

func transpileIf(b *bytes.Buffer, scp *scope, n *ast.If) {
	b.WriteString("if (")
	transpileExpr(b, scp, n.Cond)
	b.WriteString(") ")
	transpileBlock(b, scp, n.Body)

	if n.Else != nil {
		b.WriteString(" else ")
		if n.Else.Cond != nil {
			transpileIf(b, scp, n.Else)
		} else {
			transpileBlock(b, scp, n.Else.Body)
		}
	}
}

func transpileFor(b *bytes.Buffer, scp *scope, n *ast.For) {
	b.WriteString("while (")
	transpileExpr(b, scp, n.Cond)
	b.WriteString(") ")
	transpileBlock(b, scp, n.Body)
}

func transpileForRange(b *bytes.Buffer, scp *scope, n *ast.ForRange) {
	b.WriteString("for (int ")
	b.WriteString(n.Counter.Tk.Val)
	b.WriteString(" = ")
	transpileExpr(b, scp, n.Start)
	b.WriteString("; ")
	b.WriteString(n.Counter.Tk.Val)
	b.WriteString(" < ")
	transpileExpr(b, scp, n.Stop)
	b.WriteString("; ")
	b.WriteString(n.Counter.Tk.Val)
	b.WriteString(" += ")

	if n.Step != nil {
		transpileExpr(b, scp, n.Step)
	} else {
		b.WriteByte('1')
	}

	b.WriteString(") ")
	transpileBlock(b, scp, n.Body)
}

func transpileForEach(b *bytes.Buffer, scp *scope, n *ast.ForEach) {
	b.WriteString("for (auto ")
	b.WriteString(n.Holder.Tk.Val)
	b.WriteString(" : ")
	transpileExpr(b, scp, n.Iterator)
	b.WriteString(")")
	transpileBlock(b, scp, n.Body)
}

func transpileWhen(b *bytes.Buffer, scp *scope, n *ast.When) {
	b.WriteString("switch (")
	transpileExpr(b, scp, n.Val)
	b.WriteString(") {")

	for i, c := range n.Clauses {
		b.WriteString("case ")
		transpileExpr(b, scp, c)
		b.WriteByte(':')
		transpileBlock(b, scp, n.Bodies[i])
		b.WriteString("break;")
	}

	if n.Otherwise != nil {
		b.WriteString("default: ")
		transpileBlock(b, scp, n.Otherwise)
		b.WriteString("break;")
	}
	b.WriteByte('}')
}

func transpileAssign(b *bytes.Buffer, scp *scope, n *ast.Assign) {
	transpileExpr(b, scp, n.Lt)
	b.WriteByte(' ')
	b.WriteString(n.Op.Val)
	b.WriteByte(' ')
	transpileExpr(b, scp, n.Rt)
	b.WriteByte(';')
}

func transpileFlowControl(b *bytes.Buffer, _ *scope, n *ast.FlowControl) {
	b.WriteString(operators[n.Tk.Kind])
	b.WriteByte(';')
}

func transpileBlock(b *bytes.Buffer, scp *scope, n *ast.Block) {
	b.WriteByte('{')
	for _, s := range n.Stmts {
		switch s.Kind() {
		case node.Return:
			transpileReturn(b, scp, s.(*ast.Return))
		case node.StmtWrapper:
			transpileStmtWrapper(b, scp, s.(*ast.StmtWrapper))
		case node.If:
			transpileIf(b, scp, s.(*ast.If))
		case node.For:
			transpileFor(b, scp, s.(*ast.For))
		case node.ForEach:
			transpileForEach(b, scp, s.(*ast.ForEach))
		case node.ForRange:
			transpileForRange(b, scp, s.(*ast.ForRange))
		case node.When:
			transpileWhen(b, scp, s.(*ast.When))
		case node.Assign:
			transpileAssign(b, scp, s.(*ast.Assign))
		case node.FlowControl:
			transpileFlowControl(b, scp, s.(*ast.FlowControl))
		}
	}
	b.WriteByte('}')
}
