package main

type Resolver struct {
	Interp *Interpreter
	Scopes *Scopes
}

type Scopes []map[string]bool

func (s *Scopes) pop() {
	*s = (*s)[:len(*s)-1]
}

func (s *Scopes) push(scope map[string]bool) {
	*s = append(*s, scope)
}

func (r *Resolver) resolveExpr(e Expr) {
	e.Accept(r)
}

func (r *Resolver) resolveStmt(s Stmt) {
	s.Accept(r)
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
}

func (r *Resolver) visitBlockStmt(stmt *BlockStmt) (any, error) {
	return nil, nil
}

func (r *Resolver) visitVarStmt(stmt *BlockStmt) (any, error) {
	return nil, nil
}

func (i *Resolver) visitAssignExpr(expr *AssignExpr) (any, error) {
	return nil, nil
}
