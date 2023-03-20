package main

type Expr interface {
	Accept(v ExprVisitor) (any, error)
}

type ExprVisitor interface {
	VisitAssignExpr(expr *AssignExpr) (any, error)
	VisitBinaryExpr(expr *BinaryExpr) (any, error)
	VisitCallExpr(expr *CallExpr) (any, error)
	VisitGetExpr(expr *GetExpr) (any, error)
	VisitGroupingExpr(expr *GroupingExpr) (any, error)
	VisitLiteralExpr(expr *LiteralExpr) (any, error)
	VisitLogicalExpr(expr *LogicalExpr) (any, error)
	VisitSetExpr(expr *SetExpr) (any, error)
	VisitSuperExpr(expr *SuperExpr) (any, error)
	VisitThisExpr(expr *ThisExpr) (any, error)
	VisitUnaryExpr(expr *UnaryExpr) (any, error)
	VisitVarExpr(expr *VarExpr) (any, error)
}

type AssignExpr struct {
	Name  string
	Value Expr
}

func (a *AssignExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitAssignExpr(a)
}

type BinaryExpr struct {
	Left     Expr
	Right    Expr
	Operator Token
}

func (a *BinaryExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitBinaryExpr(a)
}

type CallExpr struct {
	Callee    Expr
	Paren     Token
	Arguments []Expr
}

func (a *CallExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitCallExpr(a)
}

type GetExpr struct {
	Object Expr
	Name   string
}

func (a *GetExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitGetExpr(a)
}

type GroupingExpr struct {
	Expr Expr
}

func (a *GroupingExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitGroupingExpr(a)
}

type LiteralExpr struct {
	Value any
}

func (a *LiteralExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitLiteralExpr(a)
}

type LogicalExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (a *LogicalExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitLogicalExpr(a)
}

type SetExpr struct {
}

func (a *SetExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitSetExpr(a)
}

type SuperExpr struct {
}

func (a *SuperExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitSuperExpr(a)
}

type ThisExpr struct {
}

func (a *ThisExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitThisExpr(a)
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func (a *UnaryExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitUnaryExpr(a)
}

type VarExpr struct {
	Name Token
}

func (a *VarExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitVarExpr(a)
}
