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
	Class // 23
	Else  // 24
	False // 25
	Fun   // 26
	For   // 27
	If    // 28
	Nil   // 29

	Or     // 30
	Print  // 31
	Return // 32
	Super  // 33
	This   // 34
	True   // 35
	Var    // 36
	While  // 37

	EOF // 38
)

var Keywords = map[string]TokenType{
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
