package resources

import "github.com/calummccain/coxeter/vector"

func CalcNormal(faces [][]FaceHyperbolic, p [3]float64) [3]float64 {

	h := CalcNormalEps

	return vector.Normalise3(
		vector.Sum3(vector.Scale3([3]float64{1, -1, -1}, Sdf(faces, vector.Sum3(p, vector.Scale3([3]float64{1, -1, -1}, h)))),
			vector.Sum3(vector.Scale3([3]float64{-1, -1, 1}, Sdf(faces, vector.Sum3(p, vector.Scale3([3]float64{-1, -1, 1}, h)))),
				vector.Sum3(vector.Scale3([3]float64{-1, 1, -1}, Sdf(faces, vector.Sum3(p, vector.Scale3([3]float64{-1, 1, -1}, h)))),
					vector.Scale3([3]float64{1, 1, 1}, Sdf(faces, vector.Sum3(p, vector.Scale3([3]float64{1, 1, 1}, h))))))))

}
