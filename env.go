package main

import (
	"fmt"
)

type Env struct {
	vals      map[string]any
	enclosing *Env
}

func NewEnv() *Env {
	vals := make(map[string]any, 0)
	return &Env{vals: vals}
}

func (e *Env) get(name string) (any, error) {
	if v, ok := e.vals[name]; ok {
		return v, nil
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}
	return nil, fmt.Errorf("undefined for name: %s", name)
}

func (e *Env) assign(name string, val any) error {
	if _, ok := e.vals[name]; ok {
		e.vals[name] = val
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.assign(name, val)
	}

	return fmt.Errorf("unefineable variable: %s", name)
}

func (e *Env) define(key string, val any) {
	e.vals[key] = val
}

func (e *Env) ancestor(dist int) *Env {
	environ := e
	for j := 0; j < dist; j++ {
		environ = environ.enclosing
	}
	return environ
}

func (e *Env) getAt(dist int, name string) (any, error) {
	v, ok := e.ancestor(dist).vals[name]
	if !ok {
		return nil, fmt.Errorf("can't get, at: %s, %d", name, dist)
	}
	return v, nil
}

func (e *Env) assignAt(dist int, t Token, val any) {
	e.ancestor(dist).vals[t.Lexeme] = val
}
