package transpiler

import (
	"bytes"

	"github.com/agustin-del-pino/nadia-lang/nadia/ast"
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer/token"
)

func transpileValDef(b *bytes.Buffer, scp *scope, n *ast.ValDef) {
	if n.Tk.Kind == token.Const {
		b.WriteString("const ")
	}
	transpileType(b, scp, n.Type)
	b.WriteByte(' ')
	b.WriteString(n.Name.Tk.Val)
	if n.Val != nil {
		b.WriteString(" = ")
		transpileExpr(b, scp, n.Val)
	}
	b.WriteByte(';')
}

func transpileObj(b *bytes.Buffer, scp *scope, n *ast.Obj) {
	b.WriteString("struct ")
	b.WriteString(n.Name.Tk.Val)
	b.WriteString(" {")
	for _, p := range n.Props.List {
		transpileType(b, scp, p.Type)
		b.WriteByte(' ')
		transpileIdent(b, scp, p.Name)
		b.WriteByte(';')
	}
	b.WriteString("};")
}

func transpileFunc(b *bytes.Buffer, scp *scope, n *ast.Func) {
	if n.Type != nil {
		transpileType(b, scp, n.Type)
	} else {
		b.WriteString("void")
	}

	b.WriteByte(' ')
	b.WriteString(n.Name.Tk.Val)
	b.WriteByte('(')
	for i, p := range n.Params.List {
		if i > 0 {
			b.WriteString(", ")
		}
		transpileType(b, scp, p.Type)
		b.WriteByte(' ')
		if p.Pointer {
			b.WriteByte('*')
		}
		transpileIdent(b, scp, p.Name)
	}
	b.WriteString(") ")
	transpileBlock(b, scp, n.Body)
}

func formatEventFuncName(fn *ast.EventFunc) {
	fn.Func.Name.Tk.Val = fn.Event.Tk.Val + "_" + fn.Func.Name.Tk.Val
	if fn.Tk.Kind == token.Lst {
		fn.Func.Name.Tk.Val = "lst_" + fn.Func.Name.Tk.Val
	} else {
		fn.Func.Name.Tk.Val = "ev_" + fn.Func.Name.Tk.Val
	}
}

func visit(b *bytes.Buffer, scp *scope, root *bytes.Buffer, nd *evNode, mf map[string][]*ast.EventFunc) {
	for _, ev := range nd.Children {
		transpileEventToStruct(b, scp, ev.Event)
		if nd.Init != nil {
			if nd.Init.Func != nil {
				formatEventFuncName(nd.Init)
				transpileFunc(b, scp, nd.Init.Func)
			}
			transpileEvent(b, scp, ev.Event, true)
		} else {
			transpileEvent(b, scp, ev.Event, false)
		}

		root.WriteString("if (ev_")
		root.WriteString(ev.Event.Name.Tk.Val)
		root.WriteString("_trigger()) {")

		for _, fn := range mf[ev.Event.Name.Tk.Val] {
			if fn.Func != nil {
				formatEventFuncName(fn)
				transpileFunc(b, scp, fn.Func)
				if fn.Tk.Kind == token.Lst {
					root.WriteString(fn.Func.Name.Tk.Val)
					root.WriteString("();")
				}
			} else {
				if fn.Tk.Kind == token.Lst {
					root.WriteString(fn.Link.Tk.Val)
					root.WriteString("();")
				}
			}
		}

		if len(ev.Children) > 0 {
			visit(b, scp, root, ev, mf)
		}
		root.WriteByte('}')
	}
}

func transpileSetupEvent(b *bytes.Buffer, scp *scope, tree *evTree, mf map[string][]*ast.EventFunc) {
	setup := bytes.NewBufferString("void setup() {")

	for _, l := range mf["setup"] {
		if l.Tk.Kind != token.Lst {
			continue
		}
		if l.Func != nil {
			l.Func.Name.Tk.Val = "lst_setup_" + l.Func.Name.Tk.Val
			transpileFunc(b, scp, l.Func)
			setup.WriteString(l.Func.Name.Tk.Val)
		} else {
			setup.WriteString(l.Link.Tk.Val)
		}

		setup.WriteString("();")
	}
	visit(b, scp, setup, tree.Root, mf)
	setup.WriteByte('}')
	b.Write(setup.Bytes())
}

func transpileLoopEvent(b *bytes.Buffer, scp *scope, tree *evTree, mf map[string][]*ast.EventFunc) {
	loop := bytes.NewBufferString("void loop() {")

	for _, l := range mf["loop"] {
		if l.Tk.Kind != token.Lst {
			continue
		}
		if l.Func != nil {
			l.Func.Name.Tk.Val = "lst_loop_" + l.Func.Name.Tk.Val
			transpileFunc(b, scp, l.Func)
			loop.WriteString(l.Func.Name.Tk.Val)
		} else {
			loop.WriteString(l.Link.Tk.Val)
		}
	}

	visit(b, scp, loop, tree.Root, mf)
	loop.WriteByte('}')
	b.Write(loop.Bytes())
}

func transpileEventToStruct(b *bytes.Buffer, scp *scope, n *ast.Event) {
	if len(n.Props.List) == 0 {
		return
	}
	transpileObj(b, scp, scp.Objs[n.Name.Tk.Val])
}

func transpileEvent(b *bytes.Buffer, _ *scope, n *ast.Event, init bool) {
	if len(n.Props.List) == 0 {
		return
	}
	b.WriteString(n.Name.Tk.Val)
	b.WriteString(" event_")
	b.WriteString(n.Name.Tk.Val)
	if init {
		b.WriteString(" = ev_")
		b.WriteString(n.Name.Tk.Val)
		b.WriteString("_init()")
	}
	b.WriteByte(';')
}
