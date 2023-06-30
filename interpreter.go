package main

import (
	"fmt"
	"strconv"
)

type Interpreter struct {
	env Env
}

func (i *Interpreter) interpret(stmts []Stmt) error {
	for _, v := range stmts {
		err := i.execute(v)
		if err != nil {
			return fmt.Errorf("interpreter execute: %v", err)
		}
	}
	return nil
}

func (i *Interpreter) visitAssignExpr(expr *AssignExpr) (any, error) {
	return nil, nil
}

func (i *Interpreter) visitBinaryExpr(expr *BinaryExpr) (any, error) {
	l, err := i.eval(expr.Left)
	if err != nil {
		return nil, err
	}
	r, err := i.eval(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.Type {
	case Minus:
		if lf, ok := retFloat(l); ok {
			if rf, ok := retFloat(r); ok {
				return lf - rf, nil
			}
		}
	case Plus:
		if lf, ok := retFloat(l); ok {
			if rf, ok := retFloat(r); ok {
				return lf + rf, nil
			}
		}
		if _, ok := l.(string); ok {
			if _, ok := r.(string); ok {
				return l.(string) + r.(string), nil
			}
		}
	case Slash:
		if lf, ok := retFloat(l); ok {
			if rf, ok := retFloat(r); ok {
				return lf / rf, nil
			}
		}
	case Star:
		if lf, ok := retFloat(l); ok {
			if rf, ok := retFloat(r); ok {
				return lf * rf, nil
			}
		}
	case Greater:
		if lf, ok := retFloat(l); ok {
			if rf, ok := retFloat(r); ok {
				return lf > rf, nil
			}
		}
	case GreaterEqual:
		if lf, ok := retFloat(l); ok {
			if rf, ok := retFloat(r); ok {
				return lf >= rf, nil
			}
		}
	case Less:
		if lf, ok := retFloat(l); ok {
			if rf, ok := retFloat(r); ok {
				return lf < rf, nil
			}
		}
	case LessEqual:
		if lf, ok := retFloat(l); ok {
			if rf, ok := retFloat(r); ok {
				return lf <= rf, nil
			}
		}
	case BangEqual:
		eq, err := equal(l, r)
		return !eq, err
	case EqualEqual:
		return equal(l, r)
	}
	return nil, nil
}

func (i *Interpreter) visitCallExpr(expr *CallExpr) (any, error) {
	return nil, nil
}

func (i *Interpreter) visitGetExpr(expr *GetExpr) (any, error) {
	return nil, nil
}

func (i *Interpreter) visitGroupingExpr(expr *GroupingExpr) (any, error) {
	return i.eval(expr.Expr)
}

func (i *Interpreter) visitLiteralExpr(expr *LiteralExpr) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) visitLogicalExpr(expr *LogicalExpr) (any, error) {
	return nil, nil
}

func (i *Interpreter) visitSetExpr(expr *SetExpr) (any, error) {
	return nil, nil
}

func (i *Interpreter) visitSuperExpr(expr *SuperExpr) (any, error) {
	return nil, nil
}

func (i *Interpreter) visitThisExpr(expr *ThisExpr) (any, error) {
	return nil, nil
}

func (i *Interpreter) visitUnaryExpr(expr *UnaryExpr) (any, error) {
	r, err := i.eval(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.Type {
	case Minus:
		v, ok := r.(float64)
		if !ok {
			return nil, err
		}
		return v, nil
	case Bang:
		return !truthy(r), nil
	}
	return nil, nil
}

func (i *Interpreter) visitVarExpr(expr *VarExpr) (any, error) {
	return i.env.get(expr.Name.Lexeme)
}

func (i *Interpreter) visitBlockStmt(stmt *BlockStmt) (any, error) {
	return nil, nil
}
func (i *Interpreter) visitClassStmt(stmt *ClassStmt) (any, error) {
	return nil, nil
}
func (i *Interpreter) visitExprStmt(stmt *ExprStmt) (any, error) {
	return i.eval(stmt.Expr)
}
func (i *Interpreter) visitFnStmt(stmt *FnStmt) (any, error) {
	return nil, nil
}
func (i *Interpreter) visitIfStmt(stmt *IfStmt) (any, error) {
	return nil, nil
}
func (i *Interpreter) visitPrintStmt(stmt *PrintStmt) (any, error) {
	v, err := i.eval(stmt.Expr)
	if err != nil {
		return nil, err
	}
	fmt.Println(v)
	return nil, nil
}
func (i *Interpreter) visitRetStmt(Return *RetStmt) (any, error) {
	return nil, nil
}
func (i *Interpreter) visitVarStmt(stmt *VarStmt) error {
	var v any
	var err error
	if stmt.Init != nil {
		v, err = i.eval(stmt.Init)
		if err != nil {
			return err
		}

	}
	i.env.put(stmt.Name.Lexeme, v)
	return nil
}
func (i *Interpreter) visitWhileStmt(stmt *WhileStmt) (any, error) {
	return nil, nil
}

func (i *Interpreter) eval(e Expr) (any, error) {
	return e.Accept(i)
}

func (i *Interpreter) execute(s Stmt) error {
	_, err := s.Accept(i)
	return err
}

// truthy is true for everything but false and nil.
func truthy(v any) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); !ok {
		return b
	}
	return true
}

func retFloat(v any) (float64, bool) {
	switch j := v.(type) {
	case int:
		return float64(j), true
	case float32:
		return float64(j), true
	case float64:
		return j, true
	case string:
		return strToFloat(v)
	}
	return 0, false
}

func strToFloat(v any) (float64, bool) {
	s, ok := v.(string)
	if !ok {
		return 0, false
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return f, true
}

func equal(j, k any) (bool, error) {
	if j == nil && k == nil {
		return true, nil
	}

	switch j.(type) {
	case int:
		if jf, ok := retFloat(j); ok {
			if kf, ok := retFloat(k); ok {
				return jf == kf, nil
			}
		}
	case float64:
		if kf, ok := retFloat(k); ok {
			return j == kf, nil
		}
	case string:
		if s, ok := k.(string); ok {
			return j == s, nil
		}
	}
	return false, nil
}
