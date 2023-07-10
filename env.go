package main

import (
	"fmt"
	"sync"
)

type Env struct {
	vals      sync.Map
	enclosing *Env
}

func (e *Env) get(name string) (any, error) {
	if v, ok := e.vals.Load(name); ok {
		return v, nil
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}
	return nil, fmt.Errorf("undefined for name: %s", name)
}

func (e *Env) assign(name, val any) error {
	if _, ok := e.vals.Load(name); ok {
		e.vals.Store(name, val)
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.assign(name, val)
	}

	return fmt.Errorf("unefineable variable: %s", name)
}

func (e *Env) define(key string, val any) {
	e.vals.Store(key, val)
}

func (e *Env) ancestor(dist int) *Env {
	environ := e
	for j := 0; j < dist; j++ {
		environ = environ.enclosing
	}
	return environ
}

func (e *Env) getAt(dist int, name string) (any, error) {
	v, ok := e.ancestor(dist).vals.Load(name)
	if !ok {
		return nil, fmt.Errorf("can't get, at: %s, %d", name, dist)
	}
	return v, nil
}

func (e *Env) assignAt(dist int, t Token, val any) {
	e.ancestor(dist).vals.Store(t.Lexeme, val)
}
