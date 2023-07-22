package main

type Expr interface {
	Accept(v ExprVisitor) (any, error)
}

type ExprVisitor interface {
	visitAssignExpr(expr *AssignExpr) (any, error)
	visitBinaryExpr(expr *BinaryExpr) (any, error)
	visitCallExpr(expr *CallExpr) (any, error)
	visitGetExpr(expr *GetExpr) (any, error)
	visitGroupingExpr(expr *GroupingExpr) (any, error)
	visitLiteralExpr(expr *LiteralExpr) (any, error)
	visitLogicalExpr(expr *LogicalExpr) (any, error)
	visitSetExpr(expr *SetExpr) (any, error)
	visitSuperExpr(expr *SuperExpr) (any, error)
	visitThisExpr(expr *ThisExpr) (any, error)
	visitUnaryExpr(expr *UnaryExpr) (any, error)
	visitVarExpr(expr *VarExpr) (any, error)
}

type AssignExpr struct {
	Name  string
	Value Expr
}

func (a *AssignExpr) Accept(v ExprVisitor) (any, error) {
	return v.visitAssignExpr(a)
}

type BinaryExpr struct {
	Left     Expr
	Right    Expr
	Operator Token
}

func (a *BinaryExpr) Accept(v ExprVisitor) (any, error) {
	return v.visitBinaryExpr(a)
}

type CallExpr struct {
	Callee Expr
	Paren  Token
	Args   []Expr
}

func (a *CallExpr) Accept(v ExprVisitor) (any, error) {
	return v.visitCallExpr(a)
}

type GetExpr struct {
	Object Expr
	Name   string
}

func (a *GetExpr) Accept(v ExprVisitor) (any, error) {
	return v.visitGetExpr(a)
}

type GroupingExpr struct {
	Expr Expr
}

func (a *GroupingExpr) Accept(v ExprVisitor) (any, error) {
	return v.visitGroupingExpr(a)
}

type LiteralExpr struct {
	Value any
}

func (a *LiteralExpr) Accept(v ExprVisitor) (any, error) {
	return v.visitLiteralExpr(a)
}

type LogicalExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (a *LogicalExpr) Accept(v ExprVisitor) (any, error) {
	return v.visitLogicalExpr(a)
}

type SetExpr struct {
	Object Expr
	Name   string
	Val    Expr
}

func (a *SetExpr) Accept(v ExprVisitor) (any, error) {
	return v.visitSetExpr(a)
}

type SuperExpr struct {
}

func (a *SuperExpr) Accept(v ExprVisitor) (any, error) {
	return v.visitSuperExpr(a)
}

type ThisExpr struct {
}

func (a *ThisExpr) Accept(v ExprVisitor) (any, error) {
	return v.visitThisExpr(a)
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func (a *UnaryExpr) Accept(v ExprVisitor) (any, error) {
	return v.visitUnaryExpr(a)
}

type VarExpr struct {
	Name string
}

func (a *VarExpr) Accept(v ExprVisitor) (any, error) {
	return v.visitVarExpr(a)
}
