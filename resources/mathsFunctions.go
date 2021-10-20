package resources

import "math"

func Clamp(x, min, max float64) float64 {

	return math.Min(math.Max(x, min), max)

}

func Smax(a, b float64) float64 {

	k := SmaxK

	if a+k < b {

		return b

	} else if b+k < a {

		return a

	} else {

		ab := 0.5 + 0.5*(a-b)/k

		return (1-ab)*b + ab*a + k*ab*(1-ab)

	}

}

func Smin(a, b float64) float64 {

	k := SminK

	if a+k < b {

		return a

	} else if b+k < a {

		return b

	} else {

		ab := 0.5 + 0.5*(a-b)/k

		return ab*b + (1-ab)*a - k*ab*(1-ab)

	}

}