package main

type TokenType int

const (
	STRING TokenType = iota
	NUMBER
	BOOL
	NULL
	ARRAY
	VALUE
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET
	COMMA
	COLON
)

func (t TokenType) String() string {
	switch t {
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case BOOL:
		return "BOOL"
	case NULL:
		return "NULL"
	case ARRAY:
		return "ARRAY"
	case VALUE:
		return "VALUE"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case LBRACKET:
		return "LBRACKET"
	case RBRACKET:
		return "RBRACKET"
	case COMMA:
		return "COMMA"
	case COLON:
		return "COLON"
	default:
		return "UNKNOWN"
	}
}

type Token struct {
	Type      TokenType
	Line      int
	Character int
	Value     string
}

func (t Token) String() string {
	return t.Type.String() + " " + t.Value
}
