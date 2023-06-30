package main

import "fmt"

type Env struct {
	vals map[string]any
}

func (e *Env) get(key string) (any, error) {
	v, ok := e.vals[key]
	if !ok {
		return nil, fmt.Errorf("trying to get undefined var")
	}
	return v, nil
}

func (e *Env) put(key string, val any) {
	e.vals[key] = val
}
