package resources

import (
	"math"

	"github.com/calummccain/coxeter/vector"
)

func SdfTorus(p [3]float64, ra, rb float64) float64 {

	return vector.Norm2([2]float64{vector.Norm2([2]float64{p[1], p[2]}) - ra, p[0]}) - rb

}

func SdfHyperbolic(p [3]float64, faces [][]FaceHyperbolic) float64 {

	//norm := Smin(vector.Norm3(vector.Diff3(p, [3]float64{0, 0, 0.75}))-1.0, vector.Norm3(vector.Diff3(p, [3]float64{0, 0, -0.75}))-1.0)

	//return math.Min(math.Min(vector.Norm3(vector.Diff3(p, [3]float64{2, 0, 0}))-1.0, vector.Norm3(vector.Diff3(p, [3]float64{-2, 0, 0}))-2.0), vector.Norm3(vector.Diff3(p, [3]float64{-9, 0, 0}))-4.0)

	var val, val2 float64

	norm := vector.Norm3(p) - 1.0

	for i := 0; i < len(faces); i++ {

		val2 = norm

		for j := 0; j < len(faces[i]); j++ {

			if faces[i][j].InOut {

				val2 = Smax(val2, faces[i][j].Radius-vector.Distance(p[:], faces[i][j].SphereCenter[:]))

			} else {

				val2 = Smax(val2, vector.Distance(p[:], faces[i][j].SphereCenter[:])-faces[i][j].Radius)

			}

		}

		if i == 0 {

			val = val2

		} else {

			val = math.Min(val, val2)

		}

	}

	return val

}

func SdfCube(p [3]float64, a, b, c float64) float64 {

	return Smax(math.Abs(p[0])-a, Smax(math.Abs(p[1])-b, math.Abs(p[2])-c))

	//return vector.Norm3(p) - 1
}

func SdfSphere(p [3]float64, r float64) float64 {

	return vector.Norm3(p) - r

}

func SdfSpheres(p [3]float64) float64 {

	return math.Min(math.Min(vector.Norm3(vector.Diff3(p, [3]float64{2, 0, 0}))-1.0, vector.Norm3(vector.Diff3(p, [3]float64{-2, 0, 0}))-2.0), vector.Norm3(vector.Diff3(p, [3]float64{-9, 0, 0}))-4.0)

}
