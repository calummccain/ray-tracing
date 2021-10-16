package resources

import (
	"math"

	"github.com/calummccain/coxeter/vector"
)

func RayMarch(pos, dir [3]float64, faces [][]FaceHyperbolic) ([3]float64, bool, float64) {

	var ip [3]float64
	var temp float64
	t := 0.0

	hit := false

	for i := 0; i < NumberOfIterations; i++ {

		ip = vector.Sum3(pos, vector.Scale3(dir, t))
		temp = math.Abs(Sdf(faces, ip))

		//fmt.Println(temp, i, hit)

		if temp < RayMarchEps {

			hit = true
			break

		}

		if vector.NormSquared3(ip) > 120 {

			break

		}

		t += temp

	}

	//fmt.Println(t, hit)

	return ip, hit, t
}
