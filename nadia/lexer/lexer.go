package lexer

import (
	"github.com/agustin-del-pino/nadia-lang/nadia/lexer/token"
)

type Tokenizer struct {
	src []byte
	pos int
	ln  int
	col int
}

func (t *Tokenizer) bad(tok *token.Tok, d string) {
	tok.Kind = token.Bad
	tok.Line = t.ln
	tok.Col = t.col
	tok.Val = d
}

func (t *Tokenizer) token(tok *token.Tok, s int, k token.Kind) {
	tok.Kind = k
	tok.Line = t.ln
	tok.Col = t.col
	tok.St = s
	tok.Ed = t.pos
	tok.Val = string(t.src[s:t.pos])
}

func (t *Tokenizer) Tokenize(tok *token.Tok) {
	for ; t.src[t.pos] == ' ' || t.src[t.pos] == '\n'; t.pos++ {
		if t.src[t.pos] == '\n' {
			t.ln++
			t.col = 1
		}
		if t.src[t.pos] == ' ' {
			t.col += 1
		}
	}

	if t.src[t.pos] == 0x00 {
		tok.Kind = token.EOF
		tok.Line = t.ln
		tok.Col = t.col
		tok.Val = "\x00"
		return
	}

	t.col += len(tok.Val)

	if t.src[t.pos] == '/' {
		if t.src[t.pos+1] == '/' {
			s := t.pos
			for ; !(t.src[t.pos] == 0x00 || t.src[t.pos] == '\n'); t.pos++ {
			}
			t.token(tok, s, token.Comment)
			return
		}
		if t.src[t.pos+1] == '*' {
			s := t.pos
			t.pos += 2
			for ; !(t.src[t.pos] == 0x00 || (t.src[t.pos] == '*' && t.src[t.pos+1] == '/')); t.pos++ {
			}
			if t.src[t.pos] == 0x00 {
				t.bad(tok, "comment is not close")
				return
			}
			t.pos += 2
			t.token(tok, s, token.Comment)
			return
		}
	}

	if t.src[t.pos] == '0' {
		if t.src[t.pos+1] == 'b' {
			s := t.pos
			t.pos += 2

			for ; t.src[t.pos] >= '0' && t.src[t.pos] <= '1'; t.pos++ {
			}

			if s == t.pos-2 {
				t.bad(tok, "bad formed binary number")
				return
			}

			t.token(tok, s, token.Binary)
			return
		}
		if t.src[t.pos+1] == 'x' {
			s := t.pos
			t.pos += 2

			for ; (t.src[t.pos] >= '0' && t.src[t.pos] <= '9') ||
				(t.src[t.pos] >= 'A' && t.src[t.pos] <= 'F') ||
				(t.src[t.pos] >= 'a' && t.src[t.pos] <= 'f'); t.pos++ {
			}

			if s == t.pos-2 {
				t.bad(tok, "bad formed hexadecimal number")
				return
			}

			t.token(tok, s, token.Hex)
			return
		}
		if t.src[t.pos+1] == '.' {
			s := t.pos
			t.pos += 2

			for ; t.src[t.pos] >= '0' && t.src[t.pos] <= '9'; t.pos++ {
			}

			if s == t.pos-2 {
				t.bad(tok, "bad formed float number")
				return
			}

			t.token(tok, s, token.Float)
			return
		}
		t.pos++
		t.token(tok, t.pos-1, token.Int)
		return
	}

	if t.src[t.pos] >= '1' && t.src[t.pos] <= '9' {
		s := t.pos
		for ; t.src[t.pos] >= '0' && t.src[t.pos] <= '9'; t.pos++ {
		}
		if t.src[t.pos] == '.' {
			t.pos++
			p := t.pos
			for ; t.src[t.pos] >= '0' && t.src[t.pos] <= '9'; t.pos++ {
			}

			if p == t.pos {
				t.bad(tok, "bad formed float number")
				return
			}

			t.token(tok, s, token.Float)
			return
		}
		t.token(tok, s, token.Int)
		return
	}

	if (t.src[t.pos] >= 'A' && t.src[t.pos] <= 'Z') ||
		(t.src[t.pos] >= 'a' && t.src[t.pos] <= 'z') || t.src[t.pos] == '_' {

		s := t.pos

		for ; (t.src[t.pos] >= '0' && t.src[t.pos] <= '9') ||
			(t.src[t.pos] >= 'A' && t.src[t.pos] <= 'Z') ||
			(t.src[t.pos] >= 'a' && t.src[t.pos] <= 'z') || t.src[t.pos] == '_'; t.pos++ {
		}

		t.token(tok, s, token.Ident)

		if k, ok := token.Keywords[tok.Val]; ok {
			tok.Kind = k
		}

		return
	}

	if t.src[t.pos] == '"' {
		s := t.pos
		t.pos++
		for ; !(t.src[t.pos] == 0x00 || t.src[t.pos] == '\n' || t.src[t.pos] == '"'); t.pos++ {
		}

		if t.src[t.pos] != '"' {
			t.bad(tok, "string is not close")
			return
		}
		t.pos++
		t.token(tok, s, token.String)
		return
	}

	if t.src[t.pos] == '\'' {
		s := t.pos
		t.pos++
		if t.src[t.pos] == '\\' {
			t.pos++
		}
		t.pos++

		if t.src[t.pos] != '\'' {
			t.bad(tok, "char is not close")
			return
		}

		t.pos++
		t.token(tok, s, token.Char)
		return
	}

	if t.src[t.pos] == '-' && t.src[t.pos+1] == '>' {
		t.pos += 2
		t.token(tok, t.pos-2, token.Arrow)
		return
	}

	if k, ok := token.Paired[t.src[t.pos]]; ok && t.src[t.pos+1] == '=' {
		t.pos += 2
		t.token(tok, t.pos-2, k)
		return
	}

	if k, ok := token.Unitary[t.src[t.pos]]; ok {
		t.pos++
		t.token(tok, t.pos-1, k)
		return
	}

	t.bad(tok, "unexpected char: "+string(t.src[t.pos]))
}

func sanitizes(src []byte) []byte {
	s := make([]byte, 0, len(src)+1)

	for _, b := range src {
		if b == '\r' {
			continue
		}
		if b == '\t' {
			b = ' '
		}
		s = append(s, b)
	}
	s = append(s, 0x00)
	return s[:]
}

func New(src []byte) *Tokenizer {
	return &Tokenizer{src: sanitizes(src), ln: 1, col: 1}
}
