package resources

import (
	"math"

	"github.com/calummccain/coxeter/vector"
)

func Sdf(faces [][]FaceHyperbolic, p [3]float64) float64 {

	//norm := Smin(vector.Norm3(vector.Diff3(p, [3]float64{0, 0, 0.75}))-1.0, vector.Norm3(vector.Diff3(p, [3]float64{0, 0, -0.75}))-1.0)

	// ra := 1.0
	// rb := 0.2
	// norm := vector.Norm2([2]float64{vector.Norm2([2]float64{p[1], p[2]}) - ra, p[0]}) - rb
	// return norm
	return math.Min(math.Min(vector.Norm3(vector.Diff3(p, [3]float64{2, 0, 0}))-0.5, vector.Norm3(vector.Diff3(p, [3]float64{-1.5, 0, 0}))-2), vector.Norm3(vector.Diff3(p, [3]float64{-9, 0, 0}))-4)

	// var val, val2 float64

	// norm := vector.Norm3(p) - 1.0

	// for i := 0; i < len(faces); i++ {

	// 	val2 = norm

	// 	for j := 0; j < len(faces[i]); j++ {

	// 		if faces[i][j].InOut {

	// 			val2 = Smax(val2, faces[i][j].Radius-vector.Distance(p[:], faces[i][j].SphereCenter[:]))

	// 		} else {

	// 			val2 = Smax(val2, vector.Distance(p[:], faces[i][j].SphereCenter[:])-faces[i][j].Radius)

	// 		}

	// 	}

	// 	if i == 0 {

	// 		val = val2

	// 	} else {

	// 		val = math.Min(val, val2)

	// 	}

	// }

	// return val

	//return (Smax(math.Abs(p[0]+p[1])-p[2], math.Abs(p[0]-p[1])+p[2]) - 1.0) / math.Sqrt(3.0)

	// p = [3]float64{math.Abs(p[0]), math.Abs(p[1]), math.Abs(p[2])}
	// m := p[0] + p[1] + p[2] - 2.0
	// var q [3]float64
	// if 3.0*p[0] < m {
	// 	q = p
	// } else if 3.0*p[1] < m {
	// 	q = [3]float64{p[2], p[0], p[1]}
	// } else if 3.0*p[2] < m {
	// 	q = [3]float64{p[1], p[2], p[0]}
	// } else {
	// 	return m * 0.57735027
	// }

	// k := Clamp(0.5*(q[2]-q[1]+1), 0.0, 2.0)
	// return vector.Norm3([3]float64{q[0], q[1] - 2 + k, q[2] - k})

	//return Smax(math.Abs(p[0]), Smax(math.Abs(p[1])-1, math.Abs(p[2]))) - 1

	//return vector.Norm3(p) - 1

}
