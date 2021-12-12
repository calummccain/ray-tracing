package resources

import "github.com/calummccain/coxeter/vector"

func CalcNormal(sdf func([3]float64) float64, p [3]float64) [3]float64 {

	return vector.Normalise3(
		vector.Sum3(vector.Scale3([3]float64{1, -1, -1}, sdf(vector.Sum3(p, [3]float64{CalcNormalEps, -CalcNormalEps, -CalcNormalEps}))),
			vector.Sum3(vector.Scale3([3]float64{-1, -1, 1}, sdf(vector.Sum3(p, [3]float64{-CalcNormalEps, -CalcNormalEps, CalcNormalEps}))),
				vector.Sum3(vector.Scale3([3]float64{-1, 1, -1}, sdf(vector.Sum3(p, [3]float64{-CalcNormalEps, CalcNormalEps, -CalcNormalEps}))),
					vector.Scale3([3]float64{1, 1, 1}, sdf(vector.Sum3(p, [3]float64{CalcNormalEps, CalcNormalEps, CalcNormalEps})))))))

}
