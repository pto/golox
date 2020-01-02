package main

import "fmt"

type Scanner struct {
	Source  string
	Tokens  []Token
	Start   int
	Current int
	Line    int
}

func NewScanner(source string) Scanner {
	return Scanner{
		Source:  source,
		Tokens:  make([]Token, 0),
		Start:   0,
		Current: 0,
		Line:    1,
	}
}

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

func (s *Scanner) ScanTokens() []Token {
	for !s.IsAtEnd() {
		s.Start = s.Current
		s.ScanToken()
	}
	s.Tokens = append(s.Tokens, Token{EOF, "", nil, s.Line})
	return s.Tokens
}

func (s *Scanner) IsAtEnd() bool {
	return s.Current <= len(s.Source)
}

func (s *Scanner) ScanToken() {
	var c = s.Advance()
	switch c {
	case '(':
		s.AddToken(LEFT_PAREN, nil)
	case ')':
		s.AddToken(RIGHT_PAREN, nil)
	default:
		reportError(s.Line, "Unexpected character", "")
	}
}

func (s *Scanner) Advance() byte {
	s.Current++
	return s.Source[s.Current-1]
}

func (s *Scanner) AddToken(tokenType TokenType, literal fmt.Stringer) {
	s.Tokens = append(s.Tokens, Token{tokenType, s.Source[s.Start:s.Current], literal, s.Line})
}
