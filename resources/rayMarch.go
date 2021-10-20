package resources

import (
	"math"

	"github.com/calummccain/coxeter/vector"
)

func RayMarch(sdf func([3]float64) float64, pos, dir [3]float64) ([3]float64, bool, float64) {

	var ip [3]float64
	var temp float64
	t := 0.0

	hit := false

	for i := 0; i < NumberOfIterations; i++ {

		ip = vector.Sum3(pos, vector.Scale3(dir, t))
		temp = math.Abs(sdf(ip))

		if temp < RayMarchEps {

			hit = true
			break

		}

		if vector.NormSquared3(ip) > MaxRayMarchDepth {

			break

		}

		t += temp

	}

	return ip, hit, t
}
