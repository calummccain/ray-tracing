package resources

import (
	"math"

	"github.com/calummccain/coxeter/vector"
)

//
// Union: min(sdf_a, sdf_b)
// Intersection: max(sdf_a, sdf_b)
// Subtraction: max(sdf_a, -sdf_b)
//

// Signed distance function of a torus with radii ra and rb
func SdfTorus(p [3]float64, torusConfig TorusConfig) float64 {

	return vector.Norm2([2]float64{vector.Norm2([2]float64{p[1], p[2]}) - torusConfig.TorusA, p[0]}) - torusConfig.TorusB

}

// Signed distance function of a spherical/euclidean/hyperbolic cell
func Sdf(p [3]float64, faces [][]Face, flag string) (val float64) {

	var val2 float64

	for i := 0; i < len(faces); i++ {

		if faces[i][0].Outside {

			if faces[i][0].Type == "sphere" {

				val2 = faces[i][0].Radius - vector.Distance3(p, faces[i][0].SphereCenter)

			} else {

				val2 = faces[i][0].D - vector.Dot3(p, faces[i][0].Normal)

			}

		} else {

			if faces[i][0].Type == "sphere" {

				val2 = vector.Distance3(p, faces[i][0].SphereCenter) - faces[i][0].Radius

			} else {

				val2 = -faces[i][0].D + vector.Dot3(p, faces[i][0].Normal)

			}

		}

		for j := 1; j < len(faces[i]); j++ {

			if faces[i][j].Outside {

				if faces[i][j].Type == "sphere" {

					val2 = math.Max(val2, faces[i][j].Radius-vector.Distance3(p, faces[i][j].SphereCenter))

				} else {

					val2 = math.Max(val2, faces[i][j].D-vector.Dot3(p, faces[i][j].Normal))

				}

			} else {

				if faces[i][j].Type == "sphere" {

					val2 = math.Max(val2, vector.Distance3(p, faces[i][j].SphereCenter)-faces[i][j].Radius)

				} else {

					val2 = math.Max(val2, -faces[i][j].D+vector.Dot3(p, faces[i][j].Normal))

				}

			}

		}

		if i == 0 {

			val = val2

		} else {

			val = math.Min(val, val2)

		}

	}

	if flag == "poincare" {

		val = math.Max(vector.Norm3(p)-1.0, val)

	}

	return val

}

// Signed distance function of a cuboid of lengths 2a, 2b, 2c
func SdfCube(p [3]float64, cubeConfig CubeConfig) float64 {

	return Smax(math.Abs(p[0])-cubeConfig.CubeA, Smax(math.Abs(p[1])-cubeConfig.CubeB, math.Abs(p[2])-cubeConfig.CubeC))

}

// Signed distance function of a sphere of radius r
func SdfSphere(p [3]float64, sphereConfig SphereConfig) float64 {

	return vector.Norm3(p) - sphereConfig.SphereRadius

}

// signed distance function of a row of spheres
func SdfSpheres(p [3]float64) float64 {

	return math.Min(math.Min(vector.Norm3(vector.Diff3(p, [3]float64{2, 0, 0}))-1.0, vector.Norm3(vector.Diff3(p, [3]float64{-2, 0, 0}))-2.0), vector.Norm3(vector.Diff3(p, [3]float64{-9, 0, 0}))-4.0)

}
