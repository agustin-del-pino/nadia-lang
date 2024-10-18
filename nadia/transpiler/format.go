package transpiler

import (
	"bytes"
	"strings"
)

func Format(code []byte, indent int) []byte {
	b := bytes.NewBuffer(nil)
	tab := 0
	c := append(code, 0x00)
	for i := 0; ; {
		if c[i] == 0x00 {
			break
		}

		if c[i] == '}' {
			tab--
			if tab < 0 {
				tab = 0
			}
			b.WriteByte('\n')
			b.WriteString(strings.Repeat(" ", tab*indent))
			b.WriteByte(c[i])
			if c[i+1] == ';' {
				b.WriteByte(c[i+1])
				i++
				if c[i+1] != '}' {
					b.WriteByte('\n')
				}
				i++
			} else {
				b.WriteByte('\n')
				i++
			}
			continue
		}

		if c[i] == '{' {
			tab++
			b.WriteByte(c[i])
			b.WriteByte('\n')
			b.WriteString(strings.Repeat(" ", tab*indent))
			i++
			continue
		}

		if c[i] == ';' {
			b.WriteByte(c[i])
			if c[i+1] != '}' {
				b.WriteByte('\n')
				b.WriteString(strings.Repeat(" ", tab*indent))
			}
			i++
			continue
		}

		b.WriteByte(c[i])
		i++
	}
	return b.Bytes()
}
