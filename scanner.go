package main

import (
	"strconv"
	"unicode"
	"unicode/utf8"
)

// Scanner represents the state of a Lox scanner. The positions are in bytes,
// not runes.
type Scanner struct {
	Source  string
	Tokens  []Token
	Start   int
	Current int
	Line    int
}

// NewScanner initializes a Lox scanner.
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
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"for":    For,
	"fun":    Fun,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

// ScanTokens processes the scanner source, populating and returning the
// Tokens slice.
func (s *Scanner) ScanTokens() []Token {
	for !s.IsAtEnd() {
		s.Start = s.Current
		s.ScanToken()
	}
	s.Tokens = append(s.Tokens, Token{EOF, "", nil, s.Line})
	return s.Tokens
}

// IsAtEnd indicates whether the current position is at the end of the source
// string.
func (s *Scanner) IsAtEnd() bool {
	return s.Current >= len(s.Source)
}

// ScanToken processes one token from the source string.
func (s *Scanner) ScanToken() {
	var r = s.Advance()
	switch r {
	case '(':
		s.AddToken(LeftParen, nil)
	case ')':
		s.AddToken(RightParen, nil)
	case '{':
		s.AddToken(LeftBrace, nil)
	case '}':
		s.AddToken(RightBrace, nil)
	case ',':
		s.AddToken(Comma, nil)
	case '.':
		s.AddToken(Dot, nil)
	case '-':
		s.AddToken(Minus, nil)
	case '+':
		s.AddToken(Plus, nil)
	case ';':
		s.AddToken(Semicolon, nil)
	case '*':
		s.AddToken(Star, nil)
	case '!':
		if s.Match('=') {
			s.AddToken(BangEqual, nil)
		} else {
			s.AddToken(Bang, nil)
		}
	case '=':
		if s.Match('=') {
			s.AddToken(EqualEqual, nil)
		} else {
			s.AddToken(Equal, nil)
		}
	case '<':
		if s.Match('=') {
			s.AddToken(LessEqual, nil)
		} else {
			s.AddToken(Less, nil)
		}
	case '>':
		if s.Match('=') {
			s.AddToken(GreaterEqual, nil)
		} else {
			s.AddToken(Greater, nil)
		}
	case '/':
		if s.Match('/') {
			// Comment: ignore to end of line
			for s.Peek() != '\n' && !s.IsAtEnd() {
				s.Advance()
			}
		} else {
			s.AddToken(Slash, nil)
		}
	case '\n':
		s.Line++
	case ' ', '\t', '\r':
		break
	case '"':
		s.AddString()
	default:
		if isASCIIDigit(r) {
			s.AddNumber()
		} else if unicode.IsLetter(r) || r == '_' {
			s.AddIdentifier()
		} else {
			reportError(s.Line, "Unexpected rune", "")
		}
	}
}

// Advance returns the current rune in the source string while moving the
// current position to the next rune.
func (s *Scanner) Advance() rune {
	currentRune, size := utf8.DecodeRuneInString(s.Source[s.Current:])
	s.Current += size
	return currentRune
}

// AddToken adds a token to the Tokens slice.
func (s *Scanner) AddToken(tokenType TokenType, literal interface{}) {
	s.Tokens = append(s.Tokens, Token{tokenType, s.Source[s.Start:s.Current],
		literal, s.Line})
}

// Match consumes the next rune if it matches the expected value.
func (s *Scanner) Match(expected rune) bool {
	currentRune, size := utf8.DecodeRuneInString(s.Source[s.Current:])
	if s.IsAtEnd() || currentRune != expected {
		return false
	}
	s.Current += size
	return true
}

// Peek returns the current rune without advancing the scanning position.
func (s *Scanner) Peek() rune {
	if s.IsAtEnd() {
		return '\x00'
	}
	currentRune, _ := utf8.DecodeRuneInString(s.Source[s.Current:])
	return currentRune
}

// PeekNext returns the rune after the current rune without advancing the
// scanning position.
func (s *Scanner) PeekNext() rune {
	current := s.Current
	_, size := utf8.DecodeRuneInString(s.Source[current:])
	if current+size >= len(s.Source) {
		return '\x00'
	}
	current += size
	nextRune, _ := utf8.DecodeRuneInString(s.Source[current:])
	return nextRune
}

// AddString adds a string token from the current position.
func (s *Scanner) AddString() {
	for s.Peek() != '"' && !s.IsAtEnd() {
		if s.Peek() == '\n' {
			s.Line++
		}
		s.Advance()
	}
	if s.IsAtEnd() {
		reportError(s.Line, "Unterminated string", "")
		return
	}

	// Eat closing quote
	s.Advance()

	// Save string, without the surrounding quotes (we know that a quote
	// is a single byte)
	s.AddToken(String, s.Source[s.Start+1:s.Current-1])
}

// AddNumber adds a number token from the current position.
func (s *Scanner) AddNumber() {
	for isASCIIDigit(s.Peek()) {
		s.Advance()
	}
	if s.Peek() == '.' && isASCIIDigit(s.PeekNext()) {
		s.Advance()
	}
	for isASCIIDigit(s.Peek()) {
		s.Advance()
	}
	value, err := strconv.ParseFloat(s.Source[s.Start:s.Current], 64)
	if err != nil {
		reportError(s.Line, err.Error(), "")
	}
	s.AddToken(Number, value)
}

// AddIdentifier adds a keyword or an identifier token from the current
// position.
func (s *Scanner) AddIdentifier() {
	r := s.Peek()
	for unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
		s.Advance()
		r = s.Peek()
	}
	text := s.Source[s.Start:s.Current]
	tokenType, ok := keywords[text]
	if ok {
		s.AddToken(tokenType, nil)
	} else {
		s.AddToken(Identifier, nil)
	}
}

// isASCIIDigit determines if a rune is an ASCII digit.
func isASCIIDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
