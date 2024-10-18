package token

type Kind uint8

const (
	EOF Kind = iota
	Bad
	Ident
	// Literals
	Int
	Float
	Hex
	Binary
	Char
	String
	// Keywords
	Operation
	Type
	Def
	Alias
	Var
	Const
	Func
	Obj
	Opt
	Event
	Ev
	Lst
	Include
	If
	Else
	For
	In
	Range
	When
	Is
	Otherwise
	New
	Return
	Mod
	And
	Or
	Not
	As
	Stop
	Pass
	// Delimiters
	Add
	Sub
	Mul
	Div
	Pow
	Assign
	AddAssign
	SubAssign
	MulAssign
	DivAssign
	PowAssign
	Greater
	Less
	Equal
	NotEqual
	GEqual
	LEqual
	LParen
	RParen
	LBrace
	RBrace
	LCurvy
	RCurvy
	Comma
	Dot
	Colon
	Comment
	Arrow
)

func (k Kind) String() string {
	return kindStr[k]
}

var kindStr = [...]string{
	EOF:       "\x00",
	Bad:       "",
	Ident:     "identifier",
	Int:       "integer",
	Float:     "float",
	Hex:       "hexadecimal",
	Binary:    "binary",
	Char:      "char",
	String:    "string",
	Operation: "operation",
	Type:      "type",
	Def:       "def",
	Alias:     "alias",
	Var:       "var",
	Const:     "const",
	Func:      "func",
	Obj:       "obj",
	Opt:       "ops",
	Event:     "event",
	Ev:        "ev",
	Lst:       "lst",
	Include:   "include",
	If:        "if",
	Else:      "else",
	For:       "for",
	In:        "in",
	Range:     "range",
	When:      "when",
	Is:        "is",
	Otherwise: "otherwise",
	New:       "new",
	Return:    "return",
	Mod:       "mod",
	And:       "and",
	Or:        "or",
	Not:       "not",
	As:        "as",
	Stop:      "stop",
	Pass:      "pass",
	Add:       "+",
	Sub:       "-",
	Mul:       "*",
	Div:       "/",
	Pow:       "^",
	Assign:    "=",
	AddAssign: "+=",
	SubAssign: "-=",
	MulAssign: "*=",
	DivAssign: "/=",
	PowAssign: "^=",
	Greater:   ">",
	Less:      "<",
	Equal:     "==",
	NotEqual:  "!=",
	GEqual:    ">=",
	LEqual:    "<=",
	LParen:    "(",
	RParen:    ")",
	LBrace:    "[",
	RBrace:    "]",
	LCurvy:    "{",
	RCurvy:    "}",
	Comma:     ",",
	Dot:       ".",
	Colon:     ":",
	Comment:   "//",
	Arrow:     "->",
}

var Keywords = map[string]Kind{
	"operation": Operation,
	"type":      Type,
	"def":       Def,
	"alias":     Alias,
	"var":       Var,
	"const":     Const,
	"func":      Func,
	"obj":       Obj,
	"ops":       Opt,
	"event":     Event,
	"ev":        Ev,
	"lst":       Lst,
	"include":   Include,
	"if":        If,
	"else":      Else,
	"for":       For,
	"in":        In,
	"range":     Range,
	"when":      When,
	"is":        Is,
	"otherwise": Otherwise,
	"new":       New,
	"return":    Return,
	"mod":       Mod,
	"and":       And,
	"or":        Or,
	"not":       Not,
	"as":        As,
	"pass":      Pass,
	"stop":      Stop,
}

var Unitary = map[byte]Kind{
	'+': Add,
	'-': Sub,
	'*': Mul,
	'/': Div,
	'^': Pow,
	'=': Assign,
	'>': Greater,
	'<': Less,
	'(': LParen,
	')': RParen,
	'[': LBrace,
	']': RBrace,
	'{': LCurvy,
	'}': RCurvy,
	',': Comma,
	'.': Dot,
	':': Colon,
}

var Paired = map[byte]Kind{
	'=': Equal,
	'!': NotEqual,
	'>': GEqual,
	'<': LEqual,
	'+': AddAssign,
	'-': SubAssign,
	'*': MulAssign,
	'/': DivAssign,
	'^': PowAssign,
}

type Tok struct {
	Kind Kind
	Val  string
	Line int
	Col  int
	St   int
	Ed   int
}

func (t *Tok) Clone() *Tok {
	return &Tok{
		Kind: t.Kind,
		Val:  t.Val,
		Line: t.Line,
		Col:  t.Col,
		St:   t.St,
		Ed:   t.Ed,
	}
}
