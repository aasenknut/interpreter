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
	if p.match(Fun) {
		p.curr++
		return p.funDeclaration("function")
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
		name = p.tokens[p.curr]
		p.curr++
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

func (p *Parser) funDeclaration(kind string) (Stmt, error) {
	name := p.tokens[p.curr]
	if name.Type != Identifier {
		return nil, fmt.Errorf("expected identifier: %s", kind)
	}

	p.step()
	if p.tokens[p.curr].Type != LParen {
		return nil, fmt.Errorf("expected lparen '('")
	}

	p.step()
	params := make([]Token, 0)
	for p.tokens[p.curr].Type != RParen {
		if len(params) >= 255 {
			fmt.Println("too many params in function")
		}
		ident := p.tokens[p.curr]
		if p.tokens[p.curr].Type != Identifier {
			return nil, fmt.Errorf("expected identifier - got: %v", p.tokens[p.curr].Type)
		}
		params = append(params, ident)
		p.step()
		if p.tokens[p.curr].Type == Comma {
			p.step()
		}
	}
	p.step()

	bod, err := p.block()
	if err != nil {
		return nil, fmt.Errorf("fun block: %v", err)
	}

	return &FunStmt{Name: name, Params: params, Body: bod}, nil
}

func (p *Parser) stmt() (Stmt, error) {
	if p.match(Print) {
		p.step()
		return p.printStmt()
	}
	if p.match(Return) {
		p.step()
		return p.retStmt()
	}
	if p.match(While) {
		p.step()
		return p.whileStmt()
	}
	if p.match(For) {
		p.step()
		return p.forStmt()
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

func (p *Parser) retStmt() (Stmt, error) {
	kw := p.tokens[p.curr-1]
	var val Expr
	var err error
	if p.tokens[p.curr].Type != Semicolon {
		val, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if p.tokens[p.curr].Type != Semicolon {
		return nil, fmt.Errorf("expected semicolon in ret stmt, got: %d", p.tokens[p.curr].Type)
	}
	p.step()
	return &RetStmt{
		Keyword: kw,
		Val:     val,
	}, nil
}

func (p *Parser) whileStmt() (Stmt, error) {
	if p.tokens[p.curr].Type != LParen {
		return nil, fmt.Errorf("expected lparen, while() , got: %v", p.tokens[p.curr])
	}
	cond, err := p.expression()
	if err != nil {
		return nil, fmt.Errorf("while() conditon: %v", err)
	}
	bod, err := p.stmt()
	if err != nil {
		return nil, fmt.Errorf("while() statement: %v", err)
	}
	return &WhileStmt{Cond: cond, Body: bod}, nil
}

func (p *Parser) forStmt() (Stmt, error) {
	if p.tokens[p.curr].Type != LParen {
		return nil, fmt.Errorf("for stmt - expected lparen, got: %v", p.tokens[p.curr].Type)
	}
	p.step()
	var err error
	var initialiser Stmt
	if p.match(Semicolon) {
		p.step()
		initialiser = nil
	} else if p.match(Var) {
		p.step()
		initialiser, err = p.varDeclaration()
		if err != nil {
			return nil, fmt.Errorf("for stmt, var: %v", err)
		}
	} else {
		initialiser, err = p.exprStmt()
		if err != nil {
			return nil, fmt.Errorf("for stmt: %v", err)
		}
	}
	var cond Expr = nil
	if p.tokens[p.curr].Type != Semicolon {
		cond, err = p.expression()
		if err != nil {
			return nil, fmt.Errorf("for stmt: %v", err)
		}
	}
	if p.tokens[p.curr].Type != Semicolon {
		return nil, fmt.Errorf("for stmt - expected semicolon, got: %v", p.tokens[p.curr].Type)
	}
	p.step()

	var incr Expr = nil
	if p.tokens[p.curr].Type != RParen {
		incr, err = p.expression()
		if err != nil {
			return nil, fmt.Errorf("for stmt: %v", err)
		}
	}
	if p.tokens[p.curr].Type != RParen {
		return nil, fmt.Errorf("for stmt - expected rparen, got: %v", p.tokens[p.curr].Type)
	}
	p.step()

	bod, err := p.stmt()
	if err != nil {
		return nil, fmt.Errorf("for loop statement body: %v", err)
	}

	if incr != nil {
		stmts := []Stmt{bod, &ExprStmt{incr}}
		bod = &BlockStmt{
			Stmts: stmts,
		}
	}
	if cond == nil {
		cond = &LiteralExpr{true}
	}
	bod = &WhileStmt{Cond: cond, Body: bod}

	if initialiser != nil {
		stmts := []Stmt{initialiser, bod}
		bod = &BlockStmt{Stmts: stmts}
	}
	return bod, nil
}

func (p *Parser) ifStmt() (Stmt, error) {
	p.step()
	if p.tokens[p.curr].Type != LParen {
		return nil, fmt.Errorf("if statement: expected left paren")
	}
	p.step()
	cond, err := p.expression()
	if err != nil {
		return nil, fmt.Errorf("expr from if stmt: %v", err)
	}
	if p.tokens[p.curr].Type != RParen {
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
	p.step()
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

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, fmt.Errorf("assignement: %v", err)
	}
	if p.match(Equal) {
		p.step()
		val, err := p.assignment()
		if err != nil {
			return nil, fmt.Errorf("assignement: %v", err)
		}
		if exprVal, ok := expr.(*VarExpr); ok {
			return &AssignExpr{Name: exprVal.Name, Value: val}, nil
		}
		return nil, fmt.Errorf("assignemnt - should be variable")
	}
	return expr, nil
}

func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, fmt.Errorf("or(): %v", err)
	}
	for p.match(Or) {
		p.step()
		op := p.tokens[p.curr-1]
		r, err := p.and()
		if err != nil {
			return nil, fmt.Errorf("or(): %v", err)
		}
		expr = &LogicalExpr{
			Left:     expr,
			Operator: op,
			Right:    r,
		}
	}
	return expr, nil
}

func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, fmt.Errorf("and(): %v", err)
	}
	for p.match(And) {
		op := p.tokens[p.curr]
		p.step()
		r, err := p.equality()
		if err != nil {
			return nil, fmt.Errorf("and(): %v", err)
		}
		expr = &LogicalExpr{
			Left:     expr,
			Operator: op,
			Right:    r,
		}
	}
	return expr, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
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
	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		op := p.tokens[p.curr]
		p.step()
		r, err := p.term()
		if err != nil {
			return nil, fmt.Errorf("comparison(): %v", err)
		}
		e = &BinaryExpr{Left: e, Operator: op, Right: r}
	}
	return e, err
}

func (p *Parser) term() (Expr, error) {
	e, err := p.factor()

	for p.match(Minus, Plus) {
		op := p.tokens[p.curr]
		p.step()
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
	return p.call()
}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, fmt.Errorf("call: %v", err)
	}

	for {
		if p.match(LParen) {
			p.step()
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, fmt.Errorf("call: %v", err)
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	args := make([]Expr, 0)
	for p.tokens[p.curr].Type != RParen {
		if len(args) >= 255 {
			fmt.Printf("\n[warning] more than 255 arguments")
		}
		expr, err := p.expression()
		if err != nil {
			return nil, fmt.Errorf("finish call: %v", err)
		}
		args = append(args, expr)
		if p.tokens[p.curr].Type == Comma {
			p.step()
		}
	}
	paren := p.tokens[p.curr]
	p.step()
	return &CallExpr{Callee: callee, Paren: paren, Args: args}, nil
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
		return &VarExpr{p.tokens[p.curr-1].Lexeme}, nil
	case p.match(LParen):
		p.step()
		e, err := p.expression()
		if err != nil {
			return e, err
		}
		if !p.currType(RParen) {
			return nil, fmt.Errorf("want right paren after left paren")
		}
		p.step()
		return e, nil
	}
	return nil, fmt.Errorf("get primary")
}

func (p *Parser) end() bool {
	return p.curr >= len(p.tokens)
}

func (p *Parser) currType(tp TokenType) bool {
	if p.end() {
		return false
	}
	return p.tokens[p.curr].Type == tp
}
