package resources

import "math"

func BlackBodySpectrum(l, T float64) float64 {

	l = l * 0.000000001
	// 2PIhc^2
	c1 := 0.000000000000000374183
	// hc/k_B
	c2 := 0.014388

	return c1 / (l * l * l * l * l * (math.Exp(c2/(l*T)) - 1.0))

}
