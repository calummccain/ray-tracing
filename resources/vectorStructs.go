package resources

type vector3 struct {
	x float64
	y float64
	z float64
}

func (vec1 *vector3) Sum(vec2 vector3) {
	vec1.x += vec2.x
	vec1.y += vec2.y
	vec1.z += vec2.z
}

func (vec1 *vector3) Diff(vec2 vector3) {
	vec1.x -= vec2.x
	vec1.y -= vec2.y
	vec1.z -= vec2.z
}

func (vec1 *vector3) Scale(l float64) vector3 {
	return vector3{vec1.x * l,
		vec1.y * l,
		vec1.z * l,
	}
}

func (vec1 *vector3) Dot(vec2 vector3) float64 {
	return vec1.x*vec2.x + vec1.y*vec2.y + vec1.z*vec2.z
}

func (vec1 *vector3) Norm() float64 {
	return vec1.x*vec1.x + vec1.y*vec1.y + vec1.z*vec1.z
}
