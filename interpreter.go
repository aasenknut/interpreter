package main

type Interpreter struct {
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
		return l.(float32) - r.(float32), nil
	case Plus:
		if _, ok := l.(float32); ok {
			if _, ok := r.(float32); ok {
				return l.(float32) + r.(float32), nil
			}
		}
		if _, ok := l.(string); ok {
			if _, ok := r.(string); ok {
				return l.(string) + r.(string), nil
			}
		}
	case Slash:
		return l.(float32) / r.(float32), nil
	case Star:
		return l.(float32) * r.(float32), nil
	case Greater:
		return
	case GreaterEqual:
	case Less:
	case LessEqual:
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
		v, ok := r.(float32)
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
