package resources

import (
	"math"

	"github.com/calummccain/coxeter/vector"
)

func RayMarch(sdf func([3]float64) float64, position, direction [3]float64) ([3]float64, bool, float64) {

	var tempT float64

	t := 0.0

	hit := false

	for i := 0; i < NumberOfIterations; i++ {

		tempT = math.Abs(sdf(position))

		if tempT < RayMarchEps {

			hit = true

			break

		}

		position = vector.Sum3(position, vector.Scale3(direction, tempT))
		t += tempT

		if vector.NormSquared3(position) > MaxRayMarchDepth {

			break

		}

	}

	return position, hit, t
}
