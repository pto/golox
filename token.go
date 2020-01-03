package main

import "fmt"

// Token holds a single Lox token.
type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{} // either string or float64
	Line    int
}

// String prints a debug string for a Lox token.
func (t Token) String() string {
	var literal string
	if t.Literal == nil {
		literal = "<nil>"
	} else {
		f, ok := t.Literal.(float64)
		if ok {
			literal = fmt.Sprintf("%g", f)
		} else {
			s, ok := t.Literal.(string)
			if ok {
				literal = s
			} else {
				literal = "<unknown>"
			}
		}
	}
	return fmt.Sprintf("%s %s %s %d", t.Type, t.Lexeme, literal, t.Line)
}

// TokenType indicates the Lox token type.
type TokenType int

// Enumerate the token types.
const (
	LeftParen TokenType = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual
	Identifier
	String
	Number
	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While
	EOF
)

var tokenTypeNames = [...]string{
	"LeftParen",
	"RightParen",
	"LeftBrace",
	"RightBrace",
	"Comma",
	"Dot",
	"Minus",
	"Plus",
	"Semicolon",
	"Slash",
	"Star",
	"Bang",
	"BangEqual",
	"Equal",
	"EqualEqual",
	"Greater",
	"GreaterEqual",
	"Less",
	"LessEqual",
	"Identifier",
	"String",
	"Number",
	"And",
	"Class",
	"Else",
	"False",
	"Fun",
	"For",
	"If",
	"Nil",
	"Or",
	"Print",
	"Return",
	"Super",
	"This",
	"True",
	"Var",
	"While",
	"EOF",
}

// String prints a description of the TokenType.
func (t TokenType) String() string {
	return tokenTypeNames[t]
}
