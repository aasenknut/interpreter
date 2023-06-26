package main

import (
	"fmt"
	"strconv"
	"strings"
)

type AstPrinter struct {
}

func (a *AstPrinter) Print(expr Expr) (string, error) {
	v, err := expr.Accept(a)
	if err != nil {
		return "", fmt.Errorf("print ast: %v", err)
	}
	switch s := v.(type) {
	case string:
		return s, nil
	default:
		return "", fmt.Errorf("not string")
	}
}

func (a *AstPrinter) visitAssignExpr(expr *AssignExpr) (any, error) {
	return nil, nil

}
func (a *AstPrinter) visitBinaryExpr(expr *BinaryExpr) (any, error) {
	return a.parenthesise(expr.Operator.Lexeme, expr.Left, expr.Right)

}
func (a *AstPrinter) visitCallExpr(expr *CallExpr) (any, error) {
	return nil, nil

}
func (a *AstPrinter) visitGetExpr(expr *GetExpr) (any, error) {
	return nil, nil

}
func (a *AstPrinter) visitGroupingExpr(expr *GroupingExpr) (any, error) {
	return a.parenthesise("group", expr.Expr)
}

func (a *AstPrinter) visitLiteralExpr(expr *LiteralExpr) (any, error) {
	return expr.Value, nil
}
func (a *AstPrinter) visitLogicalExpr(expr *LogicalExpr) (any, error) {
	return nil, nil

}
func (a *AstPrinter) visitSetExpr(expr *SetExpr) (any, error) {
	return nil, nil
}
func (a *AstPrinter) visitSuperExpr(expr *SuperExpr) (any, error) {
	return nil, nil
}
func (a *AstPrinter) visitThisExpr(expr *ThisExpr) (any, error) {
	return nil, nil
}
func (a *AstPrinter) visitUnaryExpr(expr *UnaryExpr) (any, error) {
	ret, err := a.parenthesise(expr.Operator.Lexeme, expr.Right)
	return ret, err
}
func (a *AstPrinter) visitVarExpr(expr *VarExpr) (any, error) {
	return nil, nil
}

func (a *AstPrinter) parenthesise(name string, exprs ...Expr) (string, error) {
	builder := strings.Builder{}
	builder.WriteString("( " + name)
	for _, j := range exprs {
		builder.WriteString(" ")
		s, err := j.Accept(a)
		if err != nil {
			return "", fmt.Errorf("trouble in paradise")
		}
		switch v := s.(type) {
		case string:
			builder.WriteString(v)
		case float64:
			str := strconv.FormatFloat(v, 'f', -1, 64)
			builder.WriteString(str)
		case int:
			str := strconv.Itoa(v)
			builder.WriteString(str)
		default:
			return "", fmt.Errorf("need string")
		}
	}
	builder.WriteString(")")
	return builder.String(), nil
}
