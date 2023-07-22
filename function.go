package main

type Function struct {
	declaration *FunStmt
	closure     *Env
}

func NewFunction() *Function {
	return &Function{}
}

func (f *Function) Call(interp *Interpreter, args []any) any {
	env := NewEnv()
	env.enclosing = CopyFrom(f.closure)
	for j := 0; j < len(f.declaration.Params); j++ {
		env.Define(f.declaration.Params[j].Lexeme, args[j])
	}
	return interp.executeBlock(f.declaration.Body, env)
}

func (f *Function) Arity() int {
	return len(f.declaration.Params)
}
