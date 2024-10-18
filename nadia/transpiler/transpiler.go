package transpiler

import (
	"bytes"
	"strconv"

	"github.com/agustin-del-pino/nadia-lang/nadia/ast"
	"github.com/agustin-del-pino/nadia-lang/nadia/ast/node"
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer/token"
)

type evNode struct {
	Event    *ast.Event
	Parent   *evNode
	HasInit  bool
	Children []*evNode
}

type evTree struct {
	Root *evNode
	Refs map[string]*evNode
}

type scope struct {
	Types map[string]*ast.TypeDef
	Defs  map[string]*ast.Def
	Objs  map[string]*ast.Obj
}

var operators = map[token.Kind]string{
	token.Not: "!",
	token.Mod: "%",
	token.And: "&&",
	token.Or:  "||",
}

func unquote(s string) string {
	return s[1 : len(s)-1]
}

func transpileType(b *bytes.Buffer, scp *scope, t *ast.Ident) {
	if ty, ok := scp.Types[t.Tk.Val]; ok {
		b.WriteString(unquote(ty.Ref.Tk.Val))
	} else if _, ok := scp.Objs[t.Tk.Val]; ok {
		b.WriteString("struct ")
		b.WriteString(t.Tk.Val)
	} else {
		b.WriteString("/* ")
		b.WriteString(t.Tk.Val)
		b.WriteString(" is not declared */")
	}
}
func transpileTypeDefault(b *bytes.Buffer, scp *scope, t *ast.Ident) {
	if ty, ok := scp.Types[t.Tk.Val]; ok {
		b.WriteString(ty.Default.Val)
	} else if _, ok := scp.Objs[t.Tk.Val]; ok {
		b.WriteString(t.Tk.Val)
		b.WriteString("{}")
	} else {
		b.WriteString("/* ")
		b.WriteString(t.Tk.Val)
		b.WriteString(" is not declared */")
	}
}

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
		}
	}
	b.WriteByte('}')
}

func transpileValDef(b *bytes.Buffer, scp *scope, n *ast.ValDef) {
	if n.Tk.Kind == token.Const {
		b.WriteString("const ")
	}
	transpileType(b, scp, n.Type)
	b.WriteByte(' ')
	b.WriteString(n.Name.Tk.Val)
	b.WriteString(" = ")
	transpileExpr(b, scp, n.Val)
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

func visit(b *bytes.Buffer, scp *scope, root *bytes.Buffer, nd *evNode, mf map[string][]*ast.EventFunc) {
	for _, ev := range nd.Children {
		transpileEventToStruct(b, scp, ev.Event)

		root.WriteString("if (ev_")
		root.WriteString(ev.Event.Name.Tk.Val)
		root.WriteString("_trigger()) {")

		for _, fn := range mf[ev.Event.Name.Tk.Val] {
			if fn.Func != nil {
				fn.Func.Name.Tk.Val = ev.Event.Name.Tk.Val + "_" + fn.Func.Name.Tk.Val

				if fn.Tk.Kind == token.Lst {
					fn.Func.Name.Tk.Val = "lst_" + fn.Func.Name.Tk.Val
				} else {
					fn.Func.Name.Tk.Val = "ev_" + fn.Func.Name.Tk.Val
				}
				transpileFunc(b, scp, fn.Func)
				if fn.Tk.Kind == token.Lst {
					root.WriteString(fn.Func.Name.Tk.Val)
					root.WriteString("(&")
					root.WriteString(ev.Event.Name.Tk.Val)
					root.WriteString(");")
				}
			} else {
				if fn.Tk.Kind == token.Lst {
					root.WriteString(fn.Link.Tk.Val)
					root.WriteString("(&")
					root.WriteString(ev.Event.Name.Tk.Val)
					root.WriteString(");")
				}
			}
		}

		transpileEvent(b, scp, ev.Event, nd.HasInit)

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

func Transpile(src *ast.Source) []byte {
	scp := scope{
		Types: map[string]*ast.TypeDef{},
		Defs:  map[string]*ast.Def{},
		Objs:  map[string]*ast.Obj{},
	}

	b := bytes.NewBuffer(nil)

	for _, t := range src.Types {
		scp.Types[t.Name.Tk.Val] = t
	}

	for _, d := range src.Defs {
		scp.Defs[d.Name.Tk.Val] = d
	}

	for _, o := range src.Objs {
		scp.Objs[o.Name.Tk.Val] = o
		transpileObj(b, &scp, o)
	}

	for _, f := range src.Funcs {
		transpileFunc(b, &scp, f)
	}

	for _, v := range src.Vals {
		transpileValDef(b, &scp, v)
	}
	setup := evTree{
		Root: &evNode{},
		Refs: make(map[string]*evNode),
	}
	loop := evTree{
		Root: &evNode{},
		Refs: make(map[string]*evNode),
	}

	for _, e := range src.Events {
		var init bool
		for _, fn := range src.EventsFunc[e.Name.Tk.Val] {
			if fn.Link != nil && fn.Link.Tk.Val == "init" {
				init = true
				break
			}
			if fn.Func != nil && fn.Func.Name.Tk.Val == "init" {
				init = true
				break
			}
		}
		scp.Objs[e.Name.Tk.Val] = &ast.Obj{
			Name:  e.Name,
			Props: e.Props,
		}
		switch e.Parent.Tk.Val {
		case "setup":
			n := evNode{
				Event:   e,
				HasInit: init,
				Parent:  setup.Root,
			}
			setup.Root.Children = append(setup.Root.Children, &n)
			setup.Refs[e.Name.Tk.Val] = &n
		case "loop":
			n := evNode{
				Event:   e,
				HasInit: init,
				Parent:  loop.Root,
			}
			loop.Root.Children = append(loop.Root.Children, &n)
			loop.Refs[e.Name.Tk.Val] = &n
		default:
			n := evNode{
				Event:   e,
				HasInit: init,
			}
			if ev, ok := setup.Refs[e.Parent.Tk.Val]; ok {
				n.Parent = ev
				ev.Children = append(ev.Children, &n)
				setup.Refs[e.Name.Tk.Val] = &n
			} else if ev, ok := loop.Refs[e.Parent.Tk.Val]; ok {
				n.Parent = ev
				ev.Children = append(ev.Children, &n)
				loop.Refs[e.Name.Tk.Val] = &n
			}
		}
	}

	transpileSetupEvent(b, &scp, &setup, src.EventsFunc)
	transpileLoopEvent(b, &scp, &loop, src.EventsFunc)

	return b.Bytes()
}
