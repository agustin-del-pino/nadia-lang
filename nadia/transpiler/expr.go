package transpiler

import (
	"bytes"
	"strconv"

	"github.com/agustin-del-pino/nadia-lang/nadia/ast"
	"github.com/agustin-del-pino/nadia-lang/nadia/ast/node"
)

func transpileIdent(b *bytes.Buffer, scp *scope, n *ast.Ident) {
	if d, ok := scp.Defs[n.Tk.Val]; ok {
		b.WriteString(unquote(d.Val.Tk.Val))
	} else {
		b.WriteString(n.Tk.Val)
	}
}

func transpileLitVal(b *bytes.Buffer, _ *scope, n *ast.LitVal) {
	b.WriteString(n.Tk.Val)
}

func transpileParenthesis(b *bytes.Buffer, scp *scope, n *ast.Parenthesis) {
	b.WriteByte('(')
	transpileExpr(b, scp, n.X)
	b.WriteByte(')')
}

func getDef(scp *scope, n ast.Node) *ast.Def {
	if n.Kind() != node.Ident {
		return nil
	}

	def, ok := scp.Defs[n.(*ast.Ident).Tk.Val]
	if !ok {
		return nil
	}

	return def
}

func transpileCall(b *bytes.Buffer, scp *scope, n *ast.Call) {
	if def := getDef(scp, n.Caller); def != nil {
		if def.Args == nil {
			b.WriteString("/*calling to a definition without arguments*/")
			return
		}

		v := append([]byte(unquote(def.Val.Tk.Val)), 0x00)
		l := len(n.Args.List)
		for i := 0; ; {
			if v[i] == 0x00 {
				break
			}
			if v[i] == '$' && (v[i+1] >= '1' && v[i+1] <= '9') {
				s := i + 1
				i += 2
				for ; v[i] >= '0' && v[i] <= '9'; i++ {
				}

				arg, _ := strconv.Atoi(string(v[s:i]))
				arg--
				if arg >= l {
					continue
				}
				transpileExpr(b, scp, n.Args.List[arg])
				continue
			}
			b.WriteByte(v[i])
			i++
		}
		return
	}
	transpileExpr(b, scp, n.Caller)
	b.WriteByte('(')
	if len(n.Args.List) > 0 {
		transpileExpr(b, scp, n.Args.List[0])
		for _, a := range n.Args.List[1:] {
			b.WriteString(", ")
			transpileExpr(b, scp, a)
		}
	}
	b.WriteByte(')')
}

func transpileSelector(b *bytes.Buffer, scp *scope, n *ast.Selector) {
	transpileExpr(b, scp, n.Path)
	b.WriteByte('.')
	transpileIdent(b, scp, n.Select)
}

func transpileUnary(b *bytes.Buffer, scp *scope, n *ast.Unary) {
	if op, ok := operators[n.Op.Kind]; ok {
		b.WriteString(op)
	} else {
		b.WriteString(n.Op.Val)
	}

	transpileExpr(b, scp, n.X)

}

func transpileBinary(b *bytes.Buffer, scp *scope, n *ast.Binary) {
	transpileExpr(b, scp, n.Lt)
	b.WriteByte(' ')
	if op, ok := operators[n.Op.Kind]; ok {
		b.WriteString(op)
	} else {
		b.WriteString(n.Op.Val)
	}
	b.WriteByte(' ')
	transpileExpr(b, scp, n.Rt)

}

func transpileNew(b *bytes.Buffer, scp *scope, n *ast.New) {
	b.WriteString("(struct ")
	b.WriteString(n.Name.Tk.Val)
	b.WriteString(") {")

	if len(n.Fields.List) > 0 {
		for i, p := range scp.Objs[n.Name.Tk.Val].Props.List {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteByte('.')
			b.WriteString(p.Name.Tk.Val)
			b.WriteString(" = ")

			var skip bool

			for _, f := range n.Fields.List {
				if f.Name.Tk.Val == p.Name.Tk.Val {
					transpileExpr(b, scp, f.Val)
					skip = true
					break
				}
			}

			if skip {
				continue
			}
			transpileTypeDefault(b, scp, p.Type)
		}
	}

	b.WriteByte('}')

}

func transpileExpr(b *bytes.Buffer, scp *scope, n ast.Expr) {
	switch n.Kind() {
	case node.Ident:
		transpileIdent(b, scp, n.(*ast.Ident))
	case node.LitVal:
		transpileLitVal(b, scp, n.(*ast.LitVal))
	case node.Parenthesis:
		transpileParenthesis(b, scp, n.(*ast.Parenthesis))
	case node.Call:
		transpileCall(b, scp, n.(*ast.Call))
	case node.Selector:
		transpileSelector(b, scp, n.(*ast.Selector))
	case node.Unary:
		transpileUnary(b, scp, n.(*ast.Unary))
	case node.Binary:
		transpileBinary(b, scp, n.(*ast.Binary))
	case node.New:
		transpileNew(b, scp, n.(*ast.New))
	}
}
