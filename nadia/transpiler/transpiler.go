package transpiler

import (
	"bytes"
	"strings"

	"github.com/agustin-del-pino/nadia-lang/nadia/ast"
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer/token"
)

type evNode struct {
	Event    *ast.Event
	Parent   *evNode
	Init     *ast.EventFunc
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
	token.Not:  "!",
	token.Mod:  "%",
	token.And:  "&&",
	token.Or:   "||",
	token.Pass: "continue",
	token.Stop: "break",
}

func unquote(s string) string {
	var b strings.Builder
	for i := 1; i < len(s)-1; i++ {
		if s[i] == '\\' {
			if s[0] == '"' {
				switch s[i+1] {
				case '\'', '"', '\\', 'n', 'r', 't', 'b', 'f', 'v', '0':
					i++
				}
			} else if s[0] == '`' && s[i+1] == '`' {
				i++
			}
		}
		b.WriteByte(s[i])
	}
	return b.String()
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
		var init *ast.EventFunc
		for _, fn := range src.EventsFunc[e.Name.Tk.Val] {
			if fn.Link != nil && fn.Link.Tk.Val == "init" {
				init = fn
				break
			}
			if fn.Func != nil && fn.Func.Name.Tk.Val == "init" {
				init = fn
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
				Event:  e,
				Init:   init,
				Parent: setup.Root,
			}
			setup.Root.Children = append(setup.Root.Children, &n)
			setup.Refs[e.Name.Tk.Val] = &n
		case "loop":
			n := evNode{
				Event:  e,
				Init:   init,
				Parent: loop.Root,
			}
			loop.Root.Children = append(loop.Root.Children, &n)
			loop.Refs[e.Name.Tk.Val] = &n
		default:
			n := evNode{
				Event: e,
				Init:  init,
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
