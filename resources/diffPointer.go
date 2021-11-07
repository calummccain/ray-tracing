package resources

func Diff3Pointer(vec1, vec2 [3]float64, out *[3]float64) {

	out[0] = vec1[0] - vec2[0]
	out[1] = vec1[1] - vec2[1]
	out[2] = vec1[2] - vec2[2]

}
