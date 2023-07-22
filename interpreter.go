package main

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
)

type Interpreter struct {
	env     *Env
	globals *Env
}

func NewInterpreter() *Interpreter {
	return &Interpreter{env: NewEnv(), globals: NewEnv()}
}

func (i *Interpreter) interpret(stmts []Stmt) error {
	for _, v := range stmts {
		_, err := i.execute(v)
		if err != nil {
			return fmt.Errorf("interpreter execute: %v", err)
		}
	}
	return nil
}

func (i *Interpreter) visitAssignExpr(expr *AssignExpr) (any, error) {
	val, err := i.eval(expr.Value)
	if err != nil {
		return nil, err
	}
	i.env.Assign(expr.Name, val)
	return val, nil
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
	callee, err := i.eval(expr.Callee)
	if err != nil {
		return nil, fmt.Errorf("visit callee, %v", err)
	}

	args := make([]any, 0)
	for _, arg := range expr.Args {
		a, err := i.eval(arg)
		if err != nil {
			return nil, fmt.Errorf("visit callee, %v", err)
		}
		args = append(args, a)
	}

	if fn, ok := callee.(Callable); ok {
		if len(args) != fn.Arity() {
			return nil, fmt.Errorf("wrong #args, expected: %d, got: %d", fn.Arity(), len(args))
		}
		return fn.Call(i, args), nil
	}
	return nil, fmt.Errorf("can not call: %v", reflect.TypeOf(callee))
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
	l, err := i.eval(expr.Left)
	if err != nil {
		return nil, fmt.Errorf("visit logical: %v", err)
	}

	if expr.Operator.Type == Or {
		if truthy(l) {
			return l, nil
		}
	} else {
		if !truthy(l) {
			return l, nil
		}
	}
	return i.eval(expr.Right)
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
	return i.env.Get(expr.Name)
}

func (i *Interpreter) visitBlockStmt(stmt *BlockStmt) (any, error) {
	e := NewEnv()
	e.enclosing = CopyFrom(i.env)
	ret := i.executeBlock(stmt.Stmts, e)
	return ret, nil
}

func (i *Interpreter) visitClassStmt(stmt *ClassStmt) (any, error) {
	return nil, nil
}
func (i *Interpreter) visitExprStmt(stmt *ExprStmt) (any, error) {
	return i.eval(stmt.Expr)
}
func (i *Interpreter) visitFunStmt(stmt *FunStmt) (any, error) {
	fun := NewFunction()
	fun.declaration = stmt
	fun.closure = i.env
	i.env.Define(stmt.Name.Lexeme, fun)
	return nil, nil
}
func (i *Interpreter) visitIfStmt(stmt *IfStmt) (any, error) {
	ok, err := i.eval(stmt.Cond)
	if err != nil {
		return nil, fmt.Errorf("visit if stmt: %v", err)
	}
	if truthy(ok) {
		return i.execute(stmt.Then)
	} else if stmt.Else != nil {
		return i.execute(stmt.Else)
	}
	return nil, nil
}
func (i *Interpreter) visitPrintStmt(stmt *PrintStmt) (any, error) {
	v, err := i.eval(stmt.Expr)
	if err != nil {
		return nil, err
	}
	if l, ok := v.(*LiteralExpr); ok {
		fmt.Println(l.Value)
	} else {
		fmt.Println(v)
	}
	return nil, nil
}
func (i *Interpreter) visitRetStmt(stmt *RetStmt) (any, error) {
	if stmt.Val != nil {
		v, err := i.eval(stmt.Val)
		if err != nil {
			return nil, fmt.Errorf("ret: %v", err)
		}
		return v, nil
	}
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
	i.env.Define(stmt.Name.Lexeme, v)
	return nil
}
func (i *Interpreter) visitWhileStmt(stmt *WhileStmt) (any, error) {
	val, err := i.eval(stmt.Cond)
	if err != nil {
		return nil, fmt.Errorf("visit whileStmt: %v", err)
	}
	for truthy(val) {
		i.execute(stmt.Body)
		val, err = i.eval(stmt.Cond)
		if err != nil {
			return nil, fmt.Errorf("visit whileStmt: %v", err)
		}
	}
	return nil, nil
}

func (i *Interpreter) eval(e Expr) (any, error) {
	return e.Accept(i)
}

func (i *Interpreter) execute(s Stmt) (any, error) {
	return s.Accept(i)
}

// truthy is true for everything but false and nil.
func truthy(v any) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		if !b {
			return b
		}
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

func (i *Interpreter) executeBlock(stmts []Stmt, e *Env) any {
	prev := CopyFrom(i.env)
	i.env = CopyFrom(e)
	var ret any
	var err error
	for _, s := range stmts {
		if retVal, ok := s.(*RetStmt); ok {
			if v, ok := retVal.Val.(*VarExpr); ok {
				ret, err = i.eval(v)
				if err != nil {
					log.Fatalf("can not get var expr %s: %v", v.Name, err)
				}
			} else {
				ret = retVal.Val
			}
		} else {
			i.execute(s)
		}
	}
	i.env = CopyFrom(prev)
	return ret
}
