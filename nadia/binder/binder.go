package binder

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/agustin-del-pino/nadia-lang/nadia/ast"
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer/token"
	"github.com/agustin-del-pino/nadia-lang/nadia/parser"
)

const (
	prx = "nad:"
)

func unquote(s string) string {
	return s[1 : len(s)-1]
}

func Bind(cwd, ext string, src *ast.Source, ign map[string]struct{}) error {
	if ign == nil {
		ign = make(map[string]struct{})
		src.Includes = append([]*ast.Include{
			{
				Path: &ast.LitVal{
					Tk: &token.Tok{
						Val: `"nad:builtin"`,
					},
				},
			},
		}, src.Includes...)
	}

	for _, inc := range src.Includes {
		if _, ok := ign[inc.Path.Tk.Val]; ok {
			continue
		}

		ign[inc.Path.Tk.Val] = struct{}{}

		pth := unquote(inc.Path.Tk.Val)

		if strings.HasPrefix(pth, prx) {
			pth = filepath.Join(ext, strings.ReplaceAll(pth, prx, ""))
		} else {
			pth = filepath.Join(cwd, pth)
		}

		pth += ".nad"

		b, bErr := os.ReadFile(pth)
		if bErr != nil {
			return bErr
		}

		s, sErr := parser.Parse(b)
		if sErr != nil {
			return sErr
		}

		bdErr := Bind(cwd, ext, s, ign)

		if bdErr != nil {
			return bdErr
		}

		src.Types = append(s.Types, src.Types...)
		src.Defs = append(s.Defs, src.Defs...)
		src.Objs = append(s.Objs, src.Objs...)
		src.Funcs = append(s.Funcs, src.Funcs...)
		src.Vals = append(s.Vals, src.Vals...)
		src.Events = append(s.Events, src.Events...)
		
		for k, v := range s.EventsFunc {
			src.EventsFunc[k] = append(v, src.EventsFunc[k]...)
		}
	}

	return nil
}
