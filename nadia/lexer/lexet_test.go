package lexer_test

import (
	"strings"
	"testing"

	"github.com/agustin-del-pino/nadia-lang/nadia/lexer"
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer/token"
	"github.com/stretchr/testify/assert"
)

func TestTokenizer_Tokenize_with_empty(t *testing.T) {
	t.Run("eof", func(t *testing.T) {
		const src = ""
		lx := setup(src)

		var tok token.Tok
		lx.Tokenize(&tok)

		assert.Equal(t, "\x00", tok.Val)
		assert.Equal(t, 1, tok.Col)
		assert.Equal(t, 1, tok.Line)
		assert.Equal(t, token.EOF, tok.Kind)
	})
	t.Run("white spaces eof", func(t *testing.T) {
		const src = "\n  \r\n \t \t"
		lx := setup(src)

		var tok token.Tok
		lx.Tokenize(&tok)

		assert.Equal(t, "\x00", tok.Val)
		assert.Equal(t, 5, tok.Col)
		assert.Equal(t, 3, tok.Line)
		assert.Equal(t, token.EOF, tok.Kind)
	})
	t.Run("after eof", func(t *testing.T) {
		const src = ""
		lx := setup(src)

		var tok token.Tok
		lx.Tokenize(&tok)
		lx.Tokenize(&tok)

		assert.Equal(t, "\x00", tok.Val)
		assert.Equal(t, 1, tok.Col)
		assert.Equal(t, 1, tok.Line)
		assert.Equal(t, token.EOF, tok.Kind)
	})
}

func TestTokenizer_Tokenize_with_valid_literals(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		const src = "0 1 2 3 4 5 6 7 8 9 1234567890"
		lx := setup(src)
		var tok token.Tok
		l := 1
		for _, ch := range strings.Split(src, " ") {
			lx.Tokenize(&tok)
			assert.Equal(t, ch, tok.Val)
			assert.Equal(t, l, tok.Col)
			assert.Equal(t, 1, tok.Line)
			assert.Equal(t, token.Int, tok.Kind)
			l += len(ch) + 1
		}
	})

	t.Run("floats", func(t *testing.T) {
		const src = "0.1 2.3 456.789 12345.67890 0.234567890"
		lx := setup(src)
		var tok token.Tok

		l := 1
		for _, ch := range strings.Split(src, " ") {
			lx.Tokenize(&tok)
			assert.Equal(t, ch, tok.Val)
			assert.Equal(t, l, tok.Col)
			assert.Equal(t, 1, tok.Line)
			assert.Equal(t, token.Float, tok.Kind)
			l += len(ch) + 1
		}
	})

	t.Run("binary", func(t *testing.T) {
		const src = "0b10101 0b11111 0b0 0b1 0b000111 0b111000"
		lx := setup(src)
		var tok token.Tok

		l := 1
		for _, ch := range strings.Split(src, " ") {
			lx.Tokenize(&tok)
			assert.Equal(t, ch, tok.Val)
			assert.Equal(t, l, tok.Col)
			assert.Equal(t, 1, tok.Line)
			assert.Equal(t, token.Binary, tok.Kind)
			l += len(ch) + 1
		}
	})

	t.Run("hex", func(t *testing.T) {
		src := make([]string, 0, 120)

		for i := '0'; i <= '9'; i++ {
			for j := 'A'; j <= 'F'; j++ {
				src = append(src, "0x"+string(i)+string(j))
			}
			for j := 'a'; j <= 'f'; j++ {
				src = append(src, "0x"+string(i)+string(j))
			}
		}

		lx := setup(strings.Join(src, " "))
		var tok token.Tok

		l := 1
		for _, ch := range src {
			lx.Tokenize(&tok)
			assert.Equal(t, ch, tok.Val)
			assert.Equal(t, l, tok.Col)
			assert.Equal(t, 1, tok.Line)
			assert.Equal(t, token.Hex, tok.Kind)
			l += len(ch) + 1
		}
	})

	t.Run("char", func(t *testing.T) {
		src := make([]string, 0, 63)

		for i := '0'; i <= '9'; i++ {
			src = append(src, "'"+string(i)+"'")
		}
		for i := 'A'; i <= 'Z'; i++ {
			src = append(src, "'"+string(i)+"'")
		}
		for i := 'a'; i <= 'z'; i++ {
			src = append(src, "'"+string(i)+"'")
		}

		src = append(src, `'\n'`)

		lx := setup(strings.Join(src, " "))
		var tok token.Tok

		l := 1
		for _, ch := range src {
			lx.Tokenize(&tok)
			assert.Equal(t, ch, tok.Val)
			assert.Equal(t, l, tok.Col)
			assert.Equal(t, 1, tok.Line)
			assert.Equal(t, token.Char, tok.Kind)
			l += len(ch) + 1
		}
	})

	t.Run("string", func(t *testing.T) {
		const src = `"1234567890qwertyuiopasdfghjklzxcvbnm_-+!#$%&/()=?¡"`
		lx := setup(src)
		var tok token.Tok

		lx.Tokenize(&tok)
		assert.Equal(t, src, tok.Val)
		assert.Equal(t, 1, tok.Col)
		assert.Equal(t, 1, tok.Line)
		assert.Equal(t, token.String, tok.Kind)
	})
}

func TestTokenizer_Tokenize_with_valid_comment(t *testing.T) {
	t.Run("single line", func(t *testing.T) {
		const src = "// this is a comment"
		lx := setup(src)
		var tok token.Tok

		lx.Tokenize(&tok)
		assert.Equal(t, src, tok.Val)
		assert.Equal(t, 1, tok.Col)
		assert.Equal(t, 1, tok.Line)
		assert.Equal(t, token.Comment, tok.Kind)
	})
	t.Run("multiline", func(t *testing.T) {
		const src = "/* this is \n a \n comment */"
		lx := setup(src)
		var tok token.Tok

		lx.Tokenize(&tok)
		assert.Equal(t, src, tok.Val)
		assert.Equal(t, 1, tok.Col)
		assert.Equal(t, 1, tok.Line)
		assert.Equal(t, token.Comment, tok.Kind)
	})
}

func TestTokenizer_Tokenize_with_keywords(t *testing.T) {
	const src = "type def alias var const func obj ops ev lst include if else for in range when is as pass stop"
	k := []token.Kind{token.Type, token.Def, token.Alias, token.Var, token.Const, token.Func, token.Obj, token.Opt, token.Ev, token.Lst, token.Include, token.If, token.Else, token.For, token.In, token.Range, token.When, token.Is, token.As, token.Pass, token.Stop}

	lx := setup(src)
	var tok token.Tok

	l := 1
	for i, ch := range strings.Split(src, " ") {
		lx.Tokenize(&tok)
		assert.Equal(t, ch, tok.Val)
		assert.Equal(t, l, tok.Col)
		assert.Equal(t, 1, tok.Line)
		assert.Equal(t, k[i], tok.Kind)
		l += len(ch) + 1
	}

}
func TestTokenizer_Tokenize_with_paired_tokens(t *testing.T) {
	const src = "== >= <= !="

	k := []token.Kind{token.Equal, token.GEqual, token.LEqual, token.NotEqual}

	lx := setup(src)
	var tok token.Tok

	l := 1
	for i, ch := range strings.Split(src, " ") {
		lx.Tokenize(&tok)
		assert.Equal(t, ch, tok.Val)
		assert.Equal(t, l, tok.Col)
		assert.Equal(t, 1, tok.Line)
		assert.Equal(t, k[i], tok.Kind)
		l += len(ch) + 1
	}
}

func TestTokenizer_Tokenize_with_single_tokens(t *testing.T) {
	const src = "+ - * / = > < ( ) [ ] { } , . :"

	k := []token.Kind{token.Add, token.Sub, token.Mul, token.Div, token.Assign, token.Greater, token.Less, token.LParen, token.RParen, token.LBrace, token.RBrace, token.LCurvy, token.RCurvy, token.Comma, token.Dot, token.Colon}

	lx := setup(src)
	var tok token.Tok

	l := 1
	for i, ch := range strings.Split(src, " ") {
		lx.Tokenize(&tok)
		assert.Equal(t, ch, tok.Val)
		assert.Equal(t, l, tok.Col)
		assert.Equal(t, 1, tok.Line)
		assert.Equal(t, k[i], tok.Kind)
		l += len(ch) + 1
	}
}

func TestTokenizer_Tokenize_with_invalid_char(t *testing.T) {
	const src = "\x1C"
	lx := setup(src)
	var tok token.Tok

	lx.Tokenize(&tok)
	assert.Equal(t, "unexpected char: \x1C", tok.Val)
	assert.Equal(t, 1, tok.Col)
	assert.Equal(t, 1, tok.Line)
	assert.Equal(t, token.Bad, tok.Kind)
}

func TestTokenizer_Tokenize_with_bad_tokens(t *testing.T) {
	t.Run("bad comment", func(t *testing.T) {
		const src = "/* asas"
		lx := setup(src)
		var tok token.Tok

		lx.Tokenize(&tok)
		assert.Equal(t, "comment is not close", tok.Val)
		assert.Equal(t, 1, tok.Col)
		assert.Equal(t, 1, tok.Line)
		assert.Equal(t, token.Bad, tok.Kind)
	})

	t.Run("bad binary", func(t *testing.T) {
		const src = "0b"
		lx := setup(src)
		var tok token.Tok

		lx.Tokenize(&tok)
		assert.Equal(t, "bad formed binary number", tok.Val)
		assert.Equal(t, 1, tok.Col)
		assert.Equal(t, 1, tok.Line)
		assert.Equal(t, token.Bad, tok.Kind)
	})

	t.Run("bad hex", func(t *testing.T) {
		const src = "0x"
		lx := setup(src)
		var tok token.Tok

		lx.Tokenize(&tok)
		assert.Equal(t, "bad formed hexadecimal number", tok.Val)
		assert.Equal(t, 1, tok.Col)
		assert.Equal(t, 1, tok.Line)
		assert.Equal(t, token.Bad, tok.Kind)
	})

	t.Run("bad zero float ", func(t *testing.T) {
		const src = "0."
		lx := setup(src)
		var tok token.Tok

		lx.Tokenize(&tok)
		assert.Equal(t, "bad formed float number", tok.Val)
		assert.Equal(t, 1, tok.Col)
		assert.Equal(t, 1, tok.Line)
		assert.Equal(t, token.Bad, tok.Kind)
	})

	t.Run("bad float ", func(t *testing.T) {
		const src = "1."
		lx := setup(src)
		var tok token.Tok

		lx.Tokenize(&tok)
		assert.Equal(t, "bad formed float number", tok.Val)
		assert.Equal(t, 1, tok.Col)
		assert.Equal(t, 1, tok.Line)
		assert.Equal(t, token.Bad, tok.Kind)
	})

	t.Run("bad string", func(t *testing.T) {
		const src = `"hello`
		lx := setup(src)
		var tok token.Tok

		lx.Tokenize(&tok)
		assert.Equal(t, "string is not close", tok.Val)
		assert.Equal(t, 1, tok.Col)
		assert.Equal(t, 1, tok.Line)
		assert.Equal(t, token.Bad, tok.Kind)
	})

	t.Run("bad char", func(t *testing.T) {
		const src = `'h`
		lx := setup(src)
		var tok token.Tok

		lx.Tokenize(&tok)
		assert.Equal(t, "char is not close", tok.Val)
		assert.Equal(t, 1, tok.Col)
		assert.Equal(t, 1, tok.Line)
		assert.Equal(t, token.Bad, tok.Kind)
	})
}

func setup(src string) *lexer.Tokenizer {
	return lexer.New([]byte(src))
}
