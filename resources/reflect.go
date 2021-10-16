package resources

import "github.com/calummccain/coxeter/vector"

func Reflect(inc, norm [3]float64) [3]float64 {

	return vector.Sum3(inc, vector.Scale3(norm, -2*vector.Dot3(norm, inc)))

}
