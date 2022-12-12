package utils

type Vector struct {
	X int
	Y int
}

func Vec(x int, y int) Vector {
	return Vector{
		X: x,
		Y: y,
	}
}

func (v *Vector) Plus(other Vector) Vector {
	return Vector{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

func (v *Vector) Minus(other Vector) Vector {
	return Vector{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}
