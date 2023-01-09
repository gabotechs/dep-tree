package ecs

import (
	"strconv"
)

type World struct {
	entities map[string]*Entity
	systems  []System
}

func NewWorld() *World {
	return &World{
		entities: map[string]*Entity{},
		systems:  []System{},
	}
}

func (w *World) WithEntity(entity *Entity) *World {
	w.entities[strconv.Itoa(len(w.entities))] = entity
	return w
}

func (w *World) WithSystem(system System) *World {
	w.systems = append(w.systems, system)
	return w
}

func (w *World) Update() error {
	for _, system := range w.systems {
		err := w.runSystem(system)
		if err != nil {
			return err
		}
	}
	return nil
}
