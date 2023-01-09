package ecs

import (
	"reflect"
)

type Entity struct {
	components map[string]any
}

func NewEntity() *Entity {
	return &Entity{
		components: make(map[string]any),
	}
}

func (e *Entity) With(component any) *Entity {
	t := reflect.TypeOf(component)
	name := t.String()
	if _, ok := e.components[name]; !ok {
		e.components[name] = component
	}
	return e
}

func (e *Entity) component(search reflect.Type) any {
	if component, ok := e.components[search.String()]; ok {
		return component
	}
	return nil
}
