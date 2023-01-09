package ecs

import (
	"reflect"
)

type System interface{}

func (w *World) runSystem(system System) error {
	typeof := reflect.TypeOf(system)
	if typeof.Kind() != reflect.Func {
		return nil
	}

	nArgs := typeof.NumIn()

	types := make([]reflect.Type, nArgs)
	for i := range types {
		types[i] = typeof.In(i)
	}
upper:
	for _, entity := range w.entities {
		components := make([]reflect.Value, nArgs)
		for i := range components {
			component := entity.component(types[i])
			if component == nil {
				continue upper
			}
			components[i] = reflect.ValueOf(component)
		}

		valueOf := reflect.ValueOf(system)
		results := valueOf.Call(components)
		if len(results) > 0 {
			if err, ok := results[0].Interface().(error); ok {
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
