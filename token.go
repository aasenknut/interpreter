package main

import "fmt"

type TokenType int

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func (t *Token) String() string {
	return fmt.Sprintf(
		"%d - %s - %v - %d", t.Type, t.Lexeme, t.Literal, t.Line)
}

const (
	// Single-character tokens.
	LParen    = iota // 0
	RParen           // 1
	LBrace           // 2
	RBrace           // 3
	Comma            // 4
	Dot              // 5
	Minus            // 6
	Plus             // 7
	Semicolon        // 8
	Slash            // 9
	Star             // 10

	// One or two character tokens.
	Bang         // 11
	BangEqual    // 12
	Equal        // 13
	EqualEqual   // 14
	Greater      // 15
	GreaterEqual // 16
	Less         // 17
	LessEqual    // 18

	// Literals.
	Identifier // 19
	String     // 20
	Number     // 21

	// Keywords.
	And   // 22
	Else  // 23
	False // 24
	Fun   // 25
	For   // 26
	If    // 27
	Nil   // 28

	Or     // 29
	Print  // 30
	Return // 31
	Super  // 32
	True   // 33
	Var    // 34
	While  // 35
)

var Keywords = map[string]TokenType{
	"and":    And,
	"else":   Else,
	"false":  False,
	"for":    For,
	"fun":    Fun,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"true":   True,
	"var":    Var,
	"while":  While,
}
