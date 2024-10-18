package node

type Kind uint8

const (
	Source Kind = iota
	TypedField
	ExprGroup
	ValueField
	OperationField
	Ident
	LitVal
	Parenthesis
	Unary
	Binary
	New
	Call
	Selector
	Stop
	Pass
	Block
	If
	For
	ForRange
	ForEach
	When
	Is
	Otherwise
	Assign
	Return
	Alias
	TypeDef
	Definition
	ValueDef
	Include
	Func
	Obj
	Opt
	EventDef
	EventFunc
	StmtWrapper
)
