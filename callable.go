package main

type Callable interface {
	Call(interp *Interpreter, args []any) any
	Arity() int
}
