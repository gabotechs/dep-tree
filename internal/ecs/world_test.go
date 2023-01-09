package ecs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWorld(t *testing.T) {
	a := require.New(t)

	type IntComponent struct {
		data int
	}

	type FloatComponent struct {
		data float32
	}

	intComponent := &IntComponent{}
	floatComponent := &FloatComponent{}

	entity := NewEntity().
		With(intComponent).
		With(floatComponent)

	system := System(func(ic *IntComponent, fc *FloatComponent) {
		ic.data++
		fc.data++
	})

	world := NewWorld().
		WithEntity(entity).
		WithSystem(system)

	err := world.Update()
	a.NoError(err)
	err = world.Update()
	a.NoError(err)

	a.Equal(2, intComponent.data)
}
