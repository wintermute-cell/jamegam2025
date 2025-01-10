package lib

type Vec2 struct {
	X float64
	Y float64
}

type Vec2I struct {
	X int
	Y int
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
func (v Vec2) Mul(s float64) Vec2 {
	return Vec2{v.X * s, v.Y * s}
}

// Div returns the division of a vector and a scalar
func (v Vec2) Div(s float64) Vec2 {
	return Vec2{v.X / s, v.Y / s}
}

// Dot returns the dot product of two vectors
func (v Vec2) Dot(v2 Vec2) float64 {
	return v.X*v2.X + v.Y*v2.Y
}

// Len returns the length of a vector
func (v Vec2) Len() float64 {
	return (v.X*v.X + v.Y*v.Y)
}

// Normalize returns the normalized vector
func (v Vec2) Normalize() Vec2 {
	l := v.Len()
	if l == 0 {
		return Vec2{0, 0}
	}
	return v.Div(v.Len())
}

// Dist returns the distance between two vectors
func (v Vec2) Dist(v2 Vec2) float64 {
	return v.Sub(v2).Len()
}

// DistSq returns the squared distance between two vectors
func (v Vec2) DistSq(v2 Vec2) float64 {
	return v.Sub(v2).Len()
}

// Lerp returns the linear interpolation between two vectors
func (v Vec2) Lerp(v2 Vec2, t float64) Vec2 {
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
	return (v.X*v.X + v.Y*v.Y)
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

// DistSq returns the squared distance between two vectors
func (v Vec2I) DistSq(v2 Vec2I) int {
	return v.Sub(v2).Len()
}

// Lerp returns the linear interpolation between two vectors
func (v Vec2I) Lerp(v2 Vec2I, t int) Vec2I {
	return v.Add(v2.Sub(v).Mul(t))
}

// ToVec2 converts a Vec2I to a Vec2
func (v Vec2I) ToVec2() Vec2 {
	return Vec2{float64(v.X), float64(v.Y)}
}

// ToVec2I converts a Vec2 to a Vec2I
func (v Vec2) ToVec2I() Vec2I {
	return Vec2I{int(v.X), int(v.Y)}
}