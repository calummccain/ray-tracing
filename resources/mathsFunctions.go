package resources

import "math"

func Clamp(x, min, max float64) float64 {

	return math.Min(math.Max(x, min), max)

}

func Smax(a, b float64) float64 {

	k := 0.01

	ab := 0.5 + 0.5*(a-b)/k

	if ab < 0 {

		return b

	} else if ab > 1 {

		return a

	} else {

		return (1-ab)*b + ab*a + k*ab*(1-ab)

	}

}

func Smin(a, b float64) float64 {

	k := 0.01

	ab := 0.5 + 0.5*(a-b)/k

	if ab < 0 {

		return a

	} else if ab > 1 {

		return b

	} else {

		return ab*b + (1-ab)*a - k*ab*(1-ab)

	}

}
