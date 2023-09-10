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

func (e *Env) Get(name string) (any, error) {
	if v, ok := e.vals[name]; ok {
		return v, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	return nil, fmt.Errorf("undefined for name: %s", name)
}

func (e *Env) Assign(name string, val any) error {
	if _, ok := e.vals[name]; ok {
		e.vals[name] = val
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(name, val)
	}

	return fmt.Errorf("unefineable variable: %s", name)
}

func (e *Env) Define(key string, val any) {
	e.vals[key] = val
}

func (e *Env) Ancestor(dist int) *Env {
	environ := e
	for j := 0; j < dist; j++ {
		environ = environ.enclosing
	}
	return environ
}

func (e *Env) GetAt(dist int, name string) (any, error) {
	v, ok := e.Ancestor(dist).vals[name]
	if !ok {
		return nil, fmt.Errorf("can't get, at: %s, %d", name, dist)
	}
	return v, nil
}

func (e *Env) AssignAt(dist int, name string, val any) {
	e.Ancestor(dist).vals[name] = val
}

func CopyFrom(src *Env) *Env {
	if src == nil {
		return &Env{}
	}
	dst := NewEnv()
	if src.enclosing != nil {
		dst.enclosing = CopyFrom(src.enclosing)
	}
	if len(src.vals) > 0 {
		dst.vals = make(map[string]any)
		for k, v := range src.vals {
			dst.vals[k] = v
		}
	}
	return dst
}
