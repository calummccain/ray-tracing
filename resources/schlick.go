package resources

import (
	"math"

	"github.com/calummccain/coxeter/vector"
)

func Schlick(inc, norm [3]float64, eta1, eta2 float64) float64 {

	//dot := vector.Dot3(inc, norm)
	// dot = math.Abs(dot)

	// r0 := (eta1 - eta2) / (eta1 + eta2)
	// r0 = r0 * r0

	// if inside {

	// 	n := eta2 / eta1
	// 	sinT := n * n * (1 - dot*dot)

	// 	if sinT > 1 {

	// 		return 1

	// 	}

	// 	dot = math.Sqrt(1 - sinT)

	// }

	// x := 1 - dot

	// return r0 + (1-r0)*x*x*x*x*x

	cos := 1 - vector.Dot3(inc, norm)

	r0 := (eta1 - eta2) / (eta1 + eta2)
	r0 = r0 * r0

	//fmt.Println(r0 + (1-r0)*cos*cos*cos*cos*cos)

	return r0 + (1-r0)*cos*cos*cos*cos*cos

}

func Fresnel(inc, norm [3]float64, eta1, eta2 float64) float64 {

	dot1 := -vector.Dot3(inc, norm)

	eta := eta1 / eta2

	if dot1 < 0 {

		dot1 = -dot1
		eta = 1.0 / eta

	}

	k := 1 - eta*eta*(1-dot1*dot1)
	var rp, rs float64

	//fmt.Println(dot1, eta, k)

	if k > 0 {

		dot2 := math.Sqrt(k)

		rs = (eta*dot1 - dot2) / (eta*dot1 + dot2)
		rp = (eta*dot2 - dot1) / (eta*dot2 + dot1)

		return 0.5 * (rs*rs + rp*rp)

	}

	return 1

}
