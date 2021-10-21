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

func SumInt(vec1, vec2 []int) []int {

	d := make([]int, 0, len(vec1))

	for i := 0; i < len(vec1); i++ {

		d = append(d, vec1[i]+vec2[i])

	}

	return d

}

func RotateX(p [3]float64, theta float64) [3]float64 {

	c := math.Cos(theta)
	s := math.Sin(theta)

	return [3]float64{p[0], c*p[1] - s*p[2], s*p[1] + c*p[2]}

}

func RotateY(p [3]float64, theta float64) [3]float64 {

	c := math.Cos(theta)
	s := math.Sin(theta)

	return [3]float64{c*p[0] - s*p[2], p[1], s*p[0] + c*p[2]}

}

func RotateZ(p [3]float64, theta float64) [3]float64 {

	c := math.Cos(theta)
	s := math.Sin(theta)

	return [3]float64{c*p[0] - s*p[1], s*p[0] + c*p[1], p[2]}

}
