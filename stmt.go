package main

type Stmt interface {
	Accept(v StmtVisitor) (any, error)
}

type StmtVisitor interface {
	visitBlockStmt(stmt *BlockStmt) (any, error)
	visitClassStmt(stmt *ClassStmt) (any, error)
	visitExprStmt(stmt *ExprStmt) (any, error)
	visitFunStmt(stmt *FunStmt) (any, error)
	visitIfStmt(stmt *IfStmt) (any, error)
	visitPrintStmt(stmt *PrintStmt) (any, error)
	visitRetStmt(Return *RetStmt) (any, error)
	visitVarStmt(stmt *VarStmt) error
	visitWhileStmt(stmt *WhileStmt) (any, error)
}

type BlockStmt struct {
	Stmts []Stmt
}

func (b *BlockStmt) Accept(v StmtVisitor) (any, error) {
	return v.visitBlockStmt(b)
}

type ClassStmt struct {
	Name       Token
	Superclass VarExpr
	Methods    []Stmt
}

func (c *ClassStmt) Accept(v StmtVisitor) (any, error) {
	return v.visitClassStmt(c)
}

type ExprStmt struct {
	Expr Expr
}

func (e *ExprStmt) Accept(v StmtVisitor) (any, error) {
	return v.visitExprStmt(e)

}

type FunStmt struct {
	Name   Token
	Params []Token
	Body   []Stmt
}

func (f *FunStmt) Accept(v StmtVisitor) (any, error) {
	return v.visitFunStmt(f)

}

type IfStmt struct {
	Cond Expr
	Then Stmt
	Else Stmt
}

func (i *IfStmt) Accept(v StmtVisitor) (any, error) {
	return v.visitIfStmt(i)
}

type PrintStmt struct {
	Expr Expr
}

func (p *PrintStmt) Accept(v StmtVisitor) (any, error) {
	return v.visitPrintStmt(p)

}

type RetStmt struct {
	Keyword Token
	Val     Expr
}

func (r *RetStmt) Accept(v StmtVisitor) (any, error) {
	return v.visitRetStmt(r)
}

type VarStmt struct {
	Name Token
	Init Expr
}

func (vs *VarStmt) Accept(v StmtVisitor) (any, error) {
	return v.visitVarStmt(vs), nil

}

type WhileStmt struct {
	Cond Expr
	Body Stmt
}

func (w *WhileStmt) Accept(v StmtVisitor) (any, error) {
	return v.visitWhileStmt(w)
}
