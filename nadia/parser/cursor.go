package parser

import (
	"fmt"

	"github.com/agustin-del-pino/nadia-lang/nadia/lexer"
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer/token"
)

type cursor token.Tok

func (c *cursor) next(t *lexer.Tokenizer) {
	t.Tokenize(c.token())
	if c.Kind == token.Bad {
		panic("bad formed token: '" + c.Val + "'")
	}
}

func (c *cursor) assert(ex string, k ...token.Kind) {
	for _, x := range k {
		if c.Kind == x {
			return
		}
	}
	panic(fmt.Sprintf("expected '%s' but got %s", ex, c.Val))
}

func (c *cursor) token() *token.Tok {
	return (*token.Tok)(c)
}
