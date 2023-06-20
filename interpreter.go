package main

type Interpreter struct {
}

func (i *Interpreter) VisitAssignExpr(expr *AssignExpr) (any, error) {
	return nil, nil
}

func (i *Interpreter) VisitBinaryExpr(expr *BinaryExpr) (any, error) {
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
	return nil, nil
}

func (i *Interpreter) VisitVarExpr(expr *VarExpr) (any, error) {
	return nil, nil
}

func (i *Interpreter) Eval(e Expr) (any, error) {
	return e.Accept(i)
}
