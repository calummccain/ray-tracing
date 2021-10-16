package resources

import (
	"math"

	"github.com/calummccain/coxeter/vector"
)

func Refract(inc, norm [3]float64, eta1, eta2 float64) ([3]float64, bool) {

	dot := -vector.Dot3(inc, norm)
	eta := eta1 / eta2

	if dot < 0 {

		dot = -dot
		norm = vector.Scale3(norm, -1.0)
		eta = 1.0 / eta

	}

	//	fmt.Println(eta)

	k := 1.0 - eta*eta*(1.0-dot*dot)

	if k <= 0.0 {

		return [3]float64{0, 0, 0}, false

	} else {

		k = math.Sqrt(k)

		//fmt.Println(vector.Norm3(inc), vector.Norm3(vector.Sum3(vector.Scale3(norm, eta*dot+k), vector.Scale3(inc, eta))))

		return vector.Normalise3(vector.Sum3(vector.Scale3(norm, -eta*dot-k), vector.Scale3(inc, eta))), true

	}

}
