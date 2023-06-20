/*
GRAMMAR
=======
expression     → equality ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
               | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")" ;
*/

package main

import "fmt"

type Parser struct {
	tokens []Token
	curr   int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() (Expr, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, fmt.Errorf("parser: %v", err)
	}
	return expr, nil
}

func (p *Parser) expression() (Expr, error) {
	eq, err := p.equality()
	return eq, err
}

func (p *Parser) equality() (Expr, error) {
	e, err := p.comparison()
	if err != nil {
		return nil, fmt.Errorf("equality(): %v", err)
	}

	for p.match(BangEqual, EqualEqual) {
		p.step()
		op := p.tokens[p.curr-1]
		right, err := p.comparison()
		if err != nil {
			return nil, fmt.Errorf("equality(): %v", err)
		}
		e = &BinaryExpr{Left: e, Operator: op, Right: right}
	}

	return e, nil
}

func (p *Parser) match(tts ...TokenType) bool {
	for _, t := range tts {
		if p.currType(t) {
			return true
		}
	}
	return false
}

func (p *Parser) step() bool {
	if p.end() {
		return false
	}
	p.curr++
	return true
}

func (p *Parser) comparison() (Expr, error) {
	e, err := p.term()
	return e, err
}

func (p *Parser) term() (Expr, error) {
	e, err := p.factor()

	for p.match(Minus, Plus) {
		p.step()
		op := p.tokens[p.curr-1]
		r, err := p.factor()
		if err != nil {
			return nil, err
		}
		e = &BinaryExpr{Left: e, Right: r, Operator: op}
	}
	return e, err
}

func (p *Parser) factor() (Expr, error) {
	e, err := p.unary()

	for p.match(Slash, Star) {
		op := p.tokens[p.curr-1]
		r, err := p.unary()
		if err != nil {
			return nil, err
		}
		e = &BinaryExpr{Left: e, Right: r, Operator: op}
	}
	return e, err
}

func (p *Parser) unary() (Expr, error) {
	if p.match(Bang, Minus) {
		op := p.tokens[p.curr]
		p.step()
		r, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{Operator: op, Right: r}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	switch {
	case p.match(False):
		p.step()
		return &LiteralExpr{Value: false}, nil
	case p.match(True):
		p.step()
		return &LiteralExpr{Value: true}, nil
	case p.match(Nil):
		p.step()
		return &LiteralExpr{Value: nil}, nil
	case p.match(Number, String):
		p.step()
		return &LiteralExpr{Value: p.tokens[p.curr-1].Literal}, nil
	case p.match(LParen):
		p.step()
		e, err := p.expression()
		if err != nil {
			return e, err
		}
		if !p.currType(RParen) {
			return nil, fmt.Errorf("want right paren after left paren")
		}
		return e, nil
	}
	return nil, fmt.Errorf("err getting primary")
}

func (p *Parser) end() bool {
	return (p.curr >= len(p.tokens)) || (p.tokens[p.curr].Type == EOF)
}

func (p *Parser) currType(tp TokenType) bool {
	if p.end() {
		return false
	}
	return p.tokens[p.curr].Type == tp
}
