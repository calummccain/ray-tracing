package resources

import (
	"math"

	hyperbolic "github.com/calummccain/coxeter/hyperbolic"
	vector "github.com/calummccain/coxeter/vector"
)

func GenerateFacesEuclidean(numFaces int, faces [][]int, localVertices []VertexEuclidean) []FaceEuclidean {

	var center4 [4]float64
	var v1, v2, v3 [3]float64
	var faceArray []FaceEuclidean

	for i := 0; i < numFaces; i++ {

		center4 = [4]float64{0, 0, 0, 0}

		for j := 0; j < len(faces[i]); j++ {

			center4 = vector.Sum4(center4, localVertices[faces[i][j]].V4)

		}

		center4 = vector.Scale4(center4, 1/vector.Norm4(center4))

		v1 = localVertices[faces[i][0]].V3
		v2 = localVertices[faces[i][1]].V3
		v3 = localVertices[faces[i][2]].V3

		faceArray = append(faceArray, FaceEuclidean{
			D:      vector.Determinant3([3][3]float64{v1, v2, v3}),
			Normal: vector.Cross3(vector.Diff3(v2, v1), vector.Diff3(v3, v1)),
		})

	}

	return faceArray

}

func GenerateFacesHyperbolic(numFaces int, faces [][]int, localVertices []VertexHyperbolic, metric byte, vv float64, model string, cellCenter [4]float64) []FaceHyperbolic {

	eps := generateEdgesHyperbolicEps

	var center4 [4]float64
	var v1, v2, v3 [4]float64
	var u1, u2, u3, centerModel [3]float64
	var faceArray []FaceHyperbolic
	var sphereCenter [3]float64
	var radius float64
	var l int
	var inOut bool
	var cellCenter3 [3]float64

	if model == "poincare" {

		cellCenter3 = hyperbolic.HyperboloidToPoincare(cellCenter)

	} else {

		cellCenter3 = hyperbolic.HyperboloidToUHP(cellCenter)

	}

	for i := 0; i < numFaces; i++ {

		center4 = [4]float64{0, 0, 0, 0}

		l = len(faces[i])

		for j := 0; j < l; j++ {

			center4 = vector.Sum4(center4, localVertices[faces[i][j]].V4)

		}

		center4 = vector.Scale4(center4, 1/math.Sqrt(math.Abs(hyperbolic.HyperbolicNorm(center4))))

		if metric == 'u' {

			v1, _ = hyperbolic.GeodesicEndpoints(localVertices[faces[i][0]].V4, localVertices[faces[i][1]].V4, vv)
			v2, _ = hyperbolic.GeodesicEndpoints(localVertices[faces[i][1]].V4, localVertices[faces[i][2]].V4, vv)
			v3, _ = hyperbolic.GeodesicEndpoints(localVertices[faces[i][2]].V4, localVertices[faces[i][3%l]].V4, vv)

		} else {

			v1 = localVertices[faces[i][0]].V4
			v2 = localVertices[faces[i][1]].V4
			v3 = localVertices[faces[i][2]].V4

		}

		if model == "uhp" {

			u1 = hyperbolic.HyperboloidToUHP(v1)
			u2 = hyperbolic.HyperboloidToUHP(v2)
			u3 = hyperbolic.HyperboloidToUHP(v3)
			centerModel = hyperbolic.HyperboloidToUHP(center4)

		} else {

			u1 = hyperbolic.HyperboloidToPoincare(v1)
			u2 = hyperbolic.HyperboloidToPoincare(v2)
			u3 = hyperbolic.HyperboloidToPoincare(v3)
			centerModel = hyperbolic.HyperboloidToPoincare(center4)

		}

		if math.Abs(vector.Determinant3([3][3]float64{vector.Diff3(u1, centerModel), vector.Diff3(u2, centerModel), vector.Diff3(u3, centerModel)})) > eps {

			sphereCenter, radius = vector.Circum4(u1, u2, u3, centerModel)

			inOut = true

			if vector.Distance(sphereCenter[:], cellCenter3[:]) < radius {

				inOut = false

			}

			faceArray = append(faceArray, FaceHyperbolic{
				Type:         "sphere",
				Radius:       radius,
				SphereCenter: sphereCenter,
				D:            0,
				Normal:       [3]float64{0, 0, 0},
				InOut:        inOut,
			})

		} else {

			if v1[2] == math.Inf(1) {

				u1 = centerModel

			} else if v2[2] == math.Inf(1) {

				u2 = centerModel

			} else if v3[2] == math.Inf(1) {

				u3 = centerModel

			}

			inOut = true

			if vector.Dot3(vector.Cross3(vector.Diff3(u2, u1), vector.Diff3(u3, u1)), cellCenter3) > 0 {

				inOut = false

			}

			faceArray = append(faceArray, FaceHyperbolic{
				Type:         "plane",
				D:            vector.Determinant3([3][3]float64{u1, u2, u3}),
				Normal:       vector.Normalise3(vector.Cross3(vector.Diff3(u2, u1), vector.Diff3(u3, u1))),
				Radius:       0,
				SphereCenter: [3]float64{0, 0, 0},
				InOut:        inOut,
			})

		}

	}

	return faceArray

}
