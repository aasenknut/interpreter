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

func (p *Parser) Parse() ([]Stmt, error) {
	stmts := make([]Stmt, 0)
	for !p.end() {
		s, err := p.declaration()
		if err != nil {
			return nil, fmt.Errorf("stmt: %v", err)
		}
		stmts = append(stmts, s)
	}
	return stmts, nil
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(Var) {
		p.curr++
		return p.varDeclaration()
	}
	s, err := p.stmt()
	if err != nil {
		return nil, fmt.Errorf("stmt declaration: %v", err)
	}
	return s, nil
}

func (p *Parser) varDeclaration() (Stmt, error) {
	var name Token
	if p.match(Identifier) {
		p.curr++
		name = p.tokens[p.curr-1]
	}
	var init Expr
	var err error
	if p.match(Equal) {
		p.curr++
		init, err = p.expression()
		if err != nil {
			return nil, fmt.Errorf("var declaration: %v", err)
		}
	}
	if p.tokens[p.curr].Type != Semicolon {
		return nil, fmt.Errorf("expected semicolon at end of var dec, got: %v", p.tokens[p.curr])
	}
	p.curr++
	return &VarStmt{Name: name, Init: init}, nil
}

func (p *Parser) stmt() (Stmt, error) {
	if p.match(Print) {
		p.step()
		return p.printStmt()
	}
	if p.tokens[p.curr].Type == LBrace {
		p.step()
		b, err := p.block()
		if err != nil {
			return nil, err
		}
		return &BlockStmt{Stmts: b}, nil
	}
	if p.match(If) {
		return p.ifStmt()
	}
	return p.exprStmt()
}

func (p *Parser) ifStmt() (Stmt, error) {
	if p.tokens[p.curr].Type != LParen {
		return nil, fmt.Errorf("if statement: expected left paren")
	}
	p.step()
	cond, err := p.expression()
	if err != nil {
		return nil, fmt.Errorf("expr from if stmt: %v", err)
	}
	p.step()
	if p.tokens[p.curr].Type != LParen {
		return nil, fmt.Errorf("if statement: expected left paren")
	}
	p.step()
	thenBranch, err := p.stmt()
	if err != nil {
		return nil, fmt.Errorf("then branch if stmt: %v", err)
	}
	var elseBranch Stmt
	elseBranch = nil
	if p.match(Else) {
		elseBranch, err = p.stmt()
		if err != nil {
			return nil, fmt.Errorf("else branch if stmt: %v", err)
		}
	}
	return &IfStmt{Cond: cond, Then: thenBranch, Else: elseBranch}, nil
}

func (p *Parser) block() ([]Stmt, error) {
	stmts := make([]Stmt, 0)
	for p.tokens[p.curr].Type != RBrace && !p.end() {
		s, err := p.declaration()
		if err != nil {
			return nil, fmt.Errorf("err declaring: %v", err)
		}
		stmts = append(stmts, s)
	}
	p.step()
	return stmts, nil
}

func (p *Parser) printStmt() (Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return nil, fmt.Errorf("print statement - expression: %v", err)
	}
	if p.tokens[p.curr].Type != Semicolon {
		return nil, fmt.Errorf("parse print - expected semicolon after print: %v", err)
	}
	p.curr++
	return &PrintStmt{val}, nil
}

func (p *Parser) exprStmt() (Stmt, error) {
	e, err := p.expression()
	if err != nil {
		return nil, fmt.Errorf("expr stmt: %v", err)
	}
	if p.tokens[p.curr].Type != Semicolon {
		return nil, fmt.Errorf("expr stmt - expected semicolon after expr: %v", err)
	}
	p.curr++
	return &ExprStmt{e}, nil
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
		return &LiteralExpr{Value: p.tokens[p.curr-1].Lexeme}, nil
	case p.match(Identifier):
		p.step()
		return &VarExpr{p.tokens[p.curr-1]}, nil
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
