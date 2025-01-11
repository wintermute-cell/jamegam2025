package lib

import "math"

type Vec2 struct {
	X float32
	Y float32
}

// NewVec2 creates a new Vec2 (float32 vector)
func NewVec2(x, y float32) Vec2 {
	return Vec2{x, y}
}

type Vec2I struct {
	X int
	Y int
}

// NewVec2I creates a new Vec2I (integer vector)
func NewVec2I(x, y int) Vec2I {
	return Vec2I{x, y}
}

// Add returns the sum of two vectors
func (v Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{v.X + v2.X, v.Y + v2.Y}
}

// Sub returns the difference of two vectors
func (v Vec2) Sub(v2 Vec2) Vec2 {
	return Vec2{v.X - v2.X, v.Y - v2.Y}
}

// Mul returns the product of a vector and a scalar
func (v Vec2) Mul(s float32) Vec2 {
	return Vec2{v.X * s, v.Y * s}
}

// Div returns the division of a vector and a scalar
func (v Vec2) Div(s float32) Vec2 {
	return Vec2{v.X / s, v.Y / s}
}

// Dot returns the dot product of two vectors
func (v Vec2) Dot(v2 Vec2) float32 {
	return v.X*v2.X + v.Y*v2.Y
}

// Len returns the length of a vector
func (v Vec2) Len() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

// Normalize returns the normalized vector
func (v Vec2) Normalize() Vec2 {
	l := v.Len()
	if l == 0 {
		return Vec2{0, 0}
	}
	return v.Div(v.Len())
}

// Rotate returns the vector rotated by an angle in degrees
func (v Vec2) Rotate(deg float32) Vec2 {
	rad := deg * math.Pi / 180
	cos := float32(math.Cos(float64(rad)))
	sin := float32(math.Sin(float64(rad)))
	return Vec2{
		X: v.X*cos - v.Y*sin,
		Y: v.X*sin + v.Y*cos,
	}
}

// Dist returns the distance between two vectors
func (v Vec2) Dist(v2 Vec2) float32 {
	return v.Sub(v2).Len()
}

// Lerp returns the linear interpolation between two vectors
func (v Vec2) Lerp(v2 Vec2, t float32) Vec2 {
	return v.Add(v2.Sub(v).Mul(t))
}

// Add returns the sum of two vectors
func (v Vec2I) Add(v2 Vec2I) Vec2I {
	return Vec2I{v.X + v2.X, v.Y + v2.Y}
}

// Sub returns the difference of two vectors
func (v Vec2I) Sub(v2 Vec2I) Vec2I {
	return Vec2I{v.X - v2.X, v.Y - v2.Y}
}

// Mul returns the product of a vector and a scalar
func (v Vec2I) Mul(s int) Vec2I {
	return Vec2I{v.X * s, v.Y * s}
}

// Div returns the division of a vector and a scalar
func (v Vec2I) Div(s int) Vec2I {
	return Vec2I{v.X / s, v.Y / s}
}

// Dot returns the dot product of two vectors
func (v Vec2I) Dot(v2 Vec2I) int {
	return v.X*v2.X + v.Y*v2.Y
}

// Len returns the length of a vector
func (v Vec2I) Len() int {
	return int(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

// Normalize returns the normalized vector
func (v Vec2I) Normalize() Vec2I {
	l := v.Len()
	if l == 0 {
		return Vec2I{0, 0}
	}
	return v.Div(v.Len())
}

// Dist returns the distance between two vectors
func (v Vec2I) Dist(v2 Vec2I) int {
	return v.Sub(v2).Len()
}

// Lerp returns the linear interpolation between two vectors
func (v Vec2I) Lerp(v2 Vec2I, t float32) Vec2I {
	// TODO: this is not that performant, could do properly...
	vf := v.ToVec2()
	v2f := v2.ToVec2()
	return vf.Lerp(v2f, t).ToVec2I()
}

// ToVec2 converts a Vec2I to a Vec2
func (v Vec2I) ToVec2() Vec2 {
	return Vec2{float32(v.X), float32(v.Y)}
}

// ToVec2I converts a Vec2 to a Vec2I
func (v Vec2) ToVec2I() Vec2I {
	return Vec2I{int(v.X), int(v.Y)}
}
