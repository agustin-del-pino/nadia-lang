package parser

import (
	"fmt"

	"github.com/agustin-del-pino/nadia-lang/nadia/ast"
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer"
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer/token"
)

type set map[token.Kind]struct{}

func (s set) Add(k token.Kind) set {
	s[k] = struct{}{}
	return s
}

func (s set) Has(k token.Kind) bool {
	_, ok := s[k]
	return ok
}

func (s set) Union(st set) set {
	for k := range st {
		s[k] = struct{}{}
	}
	return s
}

func Parse(src []byte) (*ast.Source, error) {
	return func(t *lexer.Tokenizer) (n *ast.Source, err error) {
		defer func() {
			if rcv := recover(); rcv != nil {
				err = fmt.Errorf("%s", rcv)
			}
		}()

		n = new(ast.Source)
		n.EventsFunc = make(map[string][]*ast.EventFunc)

		c := (*cursor)(&token.Tok{})
		c.next(t)

		defer func() {
			if rcv := recover(); rcv != nil {
				err = fmt.Errorf("%s at %d:%d", rcv, c.Line, c.Col)
			}
		}()

		for {
			if c.Kind == token.EOF {
				break
			}

			switch c.Kind {
			case token.Type:
				n.Types = append(n.Types, parseTypeDef(t, c))
			case token.Def:
				n.Defs = append(n.Defs, parseDef(t, c))
			case token.Var, token.Const:
				n.Vals = append(n.Vals, parseValDef(t, c))
			case token.Obj:
				n.Objs = append(n.Objs, parseObj(t, c))
			case token.Event:
				n.Events = append(n.Events, parseEvent(t, c))
			case token.Func:
				n.Funcs = append(n.Funcs, parseFunc(t, c))
			case token.Ev, token.Lst:
				fn := parseEventFunc(t, c)
				n.EventsFunc[fn.Event.Tk.Val] = append(n.EventsFunc[fn.Event.Tk.Val], fn)
			case token.Include:
				n.Includes = append(n.Includes, parseInclude(t, c))
			default:
				panic(fmt.Sprintf("invalid token: '%s'", c.Val))
			}
		}
		return n, err
	}(lexer.New(src))
}
