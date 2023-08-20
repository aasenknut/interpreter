package main

import "fmt"

type Resolver struct {
	Interp *Interpreter

	// behaves like a stack
	Scopes *Scopes
}

type Scopes []map[string]bool

func (s *Scopes) pop() {
	*s = (*s)[:len(*s)-1]
}

func (s *Scopes) push(scope map[string]bool) {
	*s = append(*s, scope)
}

func (s *Scopes) alterTop(name string, v bool) {
	topScope := (*s)[len(*s)-1]
	topScope[name] = v
}

func (s *Scopes) peek() map[string]bool {
	return (*s)[len(*s)-1]
}

func (r *Resolver) resolve(val any) {

	switch v := val.(type) {
	case Stmt:
		v.Accept(r)
	case Expr:
		v.Accept(r)
	case []Stmt:
		for _, s := range v {
			r.resolve(s)
		}
	default:
		return
	}
}

func (r *Resolver) resolveLocal(expr Expr, name string) {
	for j := len(*r.Scopes) - 1; j >= 0; j-- {
		if _, ok := (*r.Scopes)[j][name]; ok {
			r.Interp.resolve(&expr, len(*r.Scopes)-1-j)

		}
		return
	}
}

func (r *Resolver) startScope() {
	s := make(map[string]bool)
	r.Scopes.push(s)
}

func (r *Resolver) endScope() {
	r.Scopes.pop()
}

func (r *Resolver) declare(name Token) {
	if len(*r.Scopes) == 0 {
		return
	}
	r.Scopes.alterTop(name.Lexeme, false)
}

func (r *Resolver) define(name Token) {
	if len(*r.Scopes) == 0 {
		return
	}
	r.Scopes.alterTop(name.Lexeme, true)
}

func (r *Resolver) visitBlockStmt(stmt *BlockStmt) (any, error) {
	return nil, nil
}

func (r *Resolver) visitVarStmt(stmt *VarStmt) error {
	r.declare(stmt.Name)
	if stmt.Init != nil {
		r.resolve(stmt.Init)
	}
	r.define(stmt.Name)
	return nil
}

func (r *Resolver) visitVarExpr(expr *VarExpr) (any, error) {
	if len(*r.Scopes) > 0 && r.Scopes.peek()[expr.Name] == false {
		return nil, fmt.Errorf("can not read local variable in its init")
	}
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) visitWhileStmt(stmt *WhileStmt) (any, error) {
	return nil, nil
}

func (i *Resolver) visitAssignExpr(expr *AssignExpr) (any, error) {
	return nil, nil
}
func (i *Resolver) visitUnaryExpr(expr *UnaryExpr) (any, error) {
	return nil, nil
}

func (i *Resolver) visitBinaryExpr(expr *BinaryExpr) (any, error) {
	return nil, nil
}

func (i *Resolver) visitCallExpr(expr *CallExpr) (any, error) {
	return nil, nil
}

func (i *Resolver) visitGetExpr(expr *GetExpr) (any, error) {
	return nil, nil
}

func (i *Resolver) visitSetExpr(expr *SetExpr) (any, error) {
	return nil, nil
}

func (i *Resolver) visitExprStmt(expr *ExprStmt) (any, error) {
	return nil, nil
}

func (i *Resolver) visitFunStmt(expr *FunStmt) (any, error) {
	return nil, nil
}

func (i *Resolver) visitIfStmt(expr *IfStmt) (any, error) {
	return nil, nil
}

func (i *Resolver) visitPrintStmt(expr *PrintStmt) (any, error) {
	return nil, nil
}

func (i *Resolver) visitRetStmt(expr *RetStmt) (any, error) {
	return nil, nil
}

func (i *Resolver) visitLiteralExpr(expr *LiteralExpr) (any, error) {
	return nil, nil
}

func (i *Resolver) visitLogicalExpr(expr *LogicalExpr) (any, error) {
	return nil, nil
}

func (i *Resolver) visitGroupingExpr(expr *GroupingExpr) (any, error) {
	return nil, nil
}
