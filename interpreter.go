package main

import "fmt"

type Interpreter struct {
}

func (i *Interpreter) interpret(e Expr) error {
	v, err := i.Eval(e)
	if err != nil {
		return fmt.Errorf("interpreter evaluating: %v", err)
	}
	fmt.Println(v)
	return nil
}

func (i *Interpreter) VisitAssignExpr(expr *AssignExpr) (any, error) {
	return nil, nil
}

func (i *Interpreter) VisitBinaryExpr(expr *BinaryExpr) (any, error) {
	l, err := i.Eval(expr.Left)
	if err != nil {
		return nil, err
	}
	r, err := i.Eval(expr.Right)
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

func (i *Interpreter) VisitCallExpr(expr *CallExpr) (any, error) {
	return nil, nil
}

func (i *Interpreter) VisitGetExpr(expr *GetExpr) (any, error) {
	return nil, nil
}

func (i *Interpreter) VisitGroupingExpr(expr *GroupingExpr) (any, error) {
	return i.Eval(expr.Expr)
}

func (i *Interpreter) VisitLiteralExpr(expr *LiteralExpr) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitLogicalExpr(expr *LogicalExpr) (any, error) {
	return nil, nil
}

func (i *Interpreter) VisitSetExpr(expr *SetExpr) (any, error) {
	return nil, nil
}

func (i *Interpreter) VisitSuperExpr(expr *SuperExpr) (any, error) {
	return nil, nil
}

func (i *Interpreter) VisitThisExpr(expr *ThisExpr) (any, error) {
	return nil, nil
}

func (i *Interpreter) VisitUnaryExpr(expr *UnaryExpr) (any, error) {
	r, err := i.Eval(expr.Right)
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

func (i *Interpreter) VisitVarExpr(expr *VarExpr) (any, error) {
	return nil, nil
}

func (i *Interpreter) Eval(e Expr) (any, error) {
	return e.Accept(i)
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
	}
	return 0, false
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
