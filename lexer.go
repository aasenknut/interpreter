package main

import "fmt"

// Lexer consumes flat sequence of input.
// Reads character bar character to create tokens.
type Lexer struct {
	Source string
	Tokens []Token
	line   int

	// start is first character in lexeme
	start int

	// current is current character in lexeme
	current int
}

func (l *Lexer) Scan() error {
	l.current = 0
	l.start = 0
	for !l.end() {
		l.start = l.current
		err := l.readToken()
		if err != nil {
			return fmt.Errorf("read token: %v", err)
		}
	}
	return nil
}

func (l *Lexer) end() bool {
	return l.current >= len(l.Source)
}

func (l *Lexer) readToken() error {
	char := l.Source[l.current]
	switch char {
	case '(':
		l.addToken(LParen, "(", "")
		l.current++
	case ')':
		l.addToken(RParen, ")", "")
		l.current++
	case '{':
		l.addToken(LBrace, "{", "")
		l.current++
	case '}':
		l.addToken(RBrace, "}", "")
		l.current++
	case ',':
		l.addToken(Comma, ",", "")
		l.current++
	case '.':
		l.addToken(Dot, ".", "")
		l.current++
	case '-':
		l.addToken(Minus, "-", "")
		l.current++
	case '+':
		l.addToken(Plus, "+", "")
		l.current++
	case ';':
		l.addToken(Semicolon, ";", "")
		l.current++
	case '*':
		l.addToken(Star, "*", "")
		l.current++
	case '!':
		if l.lookAheadFor('=') {
			l.addToken(BangEqual, "!=", "")
		} else {
			l.addToken(Bang, "!", "")
		}
	case '=':
		if l.lookAheadFor('=') {
			l.addToken(EqualEqual, "==", "")
		} else {
			l.addToken(Equal, "=", "")
		}
		l.current++
	case '<':
		if l.lookAheadFor('=') {
			l.addToken(LessEqual, "<=", "")
		} else {
			l.addToken(Less, "<", "")
		}
	case '>':
		if l.lookAheadFor('=') {
			l.addToken(GreaterEqual, ">=", "")
		} else {
			l.addToken(Greater, ">", "")
		}
	case '/':
		if l.lookAheadFor('/') {
			for l.lookAheadFor('\n') && !l.end() {
				l.current++
			}
		} else {
			l.addToken(Slash, "/", "")
		}
	case ' ':
		l.current++
	case '\r':
		l.current++
	case '\t':
		l.current++
	case '\n':
		l.line++
		l.current++
	case '"':
		// check for string
	default: // If not numeric, then identifier.
		if isNumeric(char) {
			str, err := l.digit()
			if err != nil {
				return fmt.Errorf(
					"erronous number, line: %d, err: %v", l.line, err,
				)
			}
			l.addToken(Number, str, str)
		} else if isAlpha(char) {
			str, err := l.identifier()
			if err != nil {
				return fmt.Errorf("identifier, line: %d, err: %v", l.line, err)
			}
			l.addToken(Identifier, str, "")

		} else {
			return fmt.Errorf("unexpected character, line %d", l.line)
		}
	}
	return nil
}

func (l *Lexer) digit() (string, error) {
	for isNumeric(l.Source[l.current]) {
		l.current++
	}
	if l.Source[l.current] == '.' && isNumeric(l.Source[l.current+1]) {
		l.current++
		for isNumeric(l.Source[l.current]) {
			l.current++
		}

	}
	return l.Source[l.start:l.current], nil
}

func (l *Lexer) identifier() (string, error) {
	for isAlphaNumeric(l.Source[l.current]) {
		l.current++
	}
	return l.Source[l.start:l.current], nil
}

func (l *Lexer) lookAheadFor(want byte) bool {
	if l.current >= len(l.Source) {
		return false
	}
	return l.Source[l.current+1] == want
}

func (l *Lexer) addToken(t TokenType, lexeme, literal string) {
	l.Tokens = append(l.Tokens, Token{
		Type:    t,
		Lexeme:  lexeme,
		Literal: literal,
	})
}

func isAlphaNumeric(ch byte) bool {
	return isNumeric(ch) || isAlpha(ch)
}

func isNumeric(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
