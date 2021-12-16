package resources

import (
	"math"
)

func RayMarch(sdf func([3]float64) float64, position *vector3, direction vector3) (bool, float64) {

	var tempT float64

	t := 0.0

	hit := false

	for i := 0; i < NumberOfIterations; i++ {

		tempT = math.Abs(sdf([3]float64{position.x, position.y, position.z}))

		if tempT < RayMarchEps {

			hit = true

			break

		}

		position.Sum(direction.Scale(tempT))
		t += tempT

		if position.Norm() > MaxRayMarchDepth {

			break

		}

	}

	return hit, t
}
