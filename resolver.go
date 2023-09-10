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
			return
		}
	}
}

func (r *Resolver) resolveFun(stmt *FunStmt) {
	r.startScope()
	for _, p := range stmt.Params {
		r.declare(p)
		r.define(p)
	}
	r.resolve(stmt.Body)
	r.endScope()
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
	if len(*r.Scopes) > 0 && !r.Scopes.peek()[expr.Name] {
		return nil, fmt.Errorf("can not read local variable in its init")
	}
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) visitWhileStmt(stmt *WhileStmt) (any, error) {
	r.resolve(stmt.Cond)
	r.resolve(stmt.Body)
	return nil, nil
}

func (r *Resolver) visitAssignExpr(expr *AssignExpr) (any, error) {
	r.resolve(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}
func (r *Resolver) visitUnaryExpr(expr *UnaryExpr) (any, error) {
	r.resolve(expr.Right)
	return nil, nil
}

func (r *Resolver) visitBinaryExpr(expr *BinaryExpr) (any, error) {
	r.resolve(expr.Left)
	r.resolve(expr.Right)
	return nil, nil
}

func (r *Resolver) visitCallExpr(expr *CallExpr) (any, error) {
	r.resolve(expr.Callee)
	for _, arg := range expr.Args {
		r.resolve(arg)
	}
	return nil, nil
}

func (r *Resolver) visitGetExpr(expr *GetExpr) (any, error) {
	return nil, nil
}

func (r *Resolver) visitSetExpr(expr *SetExpr) (any, error) {
	return nil, nil
}

func (r *Resolver) visitExprStmt(stmt *ExprStmt) (any, error) {
	r.resolve(stmt.Expr)
	return nil, nil
}

func (r *Resolver) visitFunStmt(stmt *FunStmt) (any, error) {
	r.declare(stmt.Name)
	r.define(stmt.Name)
	r.resolveFun(stmt)
	return nil, nil
}

func (r *Resolver) visitIfStmt(stmt *IfStmt) (any, error) {
	r.resolve(stmt.Cond)
	r.resolve(stmt.Then)
	if stmt.Else != nil {
		r.resolve(stmt.Else)
	}
	return nil, nil
}

func (r *Resolver) visitPrintStmt(stmt *PrintStmt) (any, error) {
	r.resolve(stmt.Expr)
	return nil, nil
}

func (r *Resolver) visitRetStmt(stmt *RetStmt) (any, error) {
	if stmt.Val != nil {
		r.resolve(stmt.Val)
	}
	return nil, nil
}

func (r *Resolver) visitLiteralExpr(expr *LiteralExpr) (any, error) {
	return nil, nil
}

func (r *Resolver) visitLogicalExpr(expr *LogicalExpr) (any, error) {
	r.resolve(expr.Right)
	return nil, nil
}

func (r *Resolver) visitGroupingExpr(expr *GroupingExpr) (any, error) {
	r.resolve(expr.Expr)
	return nil, nil
}
