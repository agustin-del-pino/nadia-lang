package transpiler

import (
	"bytes"
	"strings"
)

func Format(code []byte, indent int) []byte {
	b := bytes.NewBuffer(nil)
	tab := 0
	c := append(code, 0x00)
	var prn bool
	for i := 0; c[i] != 0x00; i++ {
		if i > 0 {
			if (c[i-1] == '{' && c[i] != '}') ||
				(c[i-1] == '}' && c[i] != ';') ||
				(c[i-1] == ';' && !prn) {
				b.WriteByte('\n')
			}
			if (c[i-1] == '{' && c[i] != '}') ||
				(c[i-1] == ';' && !prn) ||
				(c[i-1] == '}' && c[i] != '}') {
				b.WriteString(strings.Repeat(" ", tab*indent))
			}
		}

		b.WriteByte(c[i])
		switch c[i] {
		case '(':
			prn = true
		case ')':
			prn = false
		case '{':
			tab++
		}
		if c[i+1] == '}' {
			tab--
		}
	}
	return b.Bytes()
}
