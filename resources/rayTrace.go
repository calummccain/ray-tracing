package resources

import (
	"math"

	"github.com/calummccain/coxeter/vector"
)

type ray struct {
	pos     [3]float64
	dir     [3]float64
	weight  float64
	rayType string
	parent  int
	inside  bool
	layer   int
}

func RayTrace(sdf func([3]float64) float64, dir, pos [3]float64, eta1 float64, eta2 [3]float64, iterations int, faces [][]FaceHyperbolic, up, left [3]float64, invHeight, invWidth float64, raysPerPixel int) ([3]float64, int, int, int, []int) {

	pixelColor := [3]float64{0, 0, 0}
	numberOfRays := 0
	numberOfMarches := 0
	depthStatistics := []int{0}

	var wShift, hShift float64
	var shiftedDir [3]float64

	var rays [][]ray
	var newRays []ray

	var newPos [3]float64
	var hit bool
	var newRayPos [3]float64

	var norm [3]float64
	var schlick float64

	var refractDir [3]float64
	var refract bool

	testRefract := true
	testReflect := true

	shift := RayTraceShift

	for l := 0; l < iterations; l++ {
		depthStatistics = append(depthStatistics, 0)
	}

	for w := 0; w < raysPerPixel; w++ {

		wShift = (float64(w) - float64(raysPerPixel-1)*0.5) * invWidth / float64(raysPerPixel)

		for h := 0; h < raysPerPixel; h++ {

			hShift = (float64(h) - float64(raysPerPixel-1)*0.5) * invHeight / float64(raysPerPixel)

			shiftedDir = vector.Normalise3(vector.Sum3(vector.Sum3(dir, vector.Scale3(up, hShift)), vector.Scale3(left, wShift)))

			for j := 0; j < 3; j += 1 {

				numberOfRays += 1

				rays = [][]ray{
					{
						{pos: pos, dir: shiftedDir, weight: 1, rayType: "initial", parent: 0, inside: false, layer: -1},
					},
				}

				k := 0
				for k < iterations && len(rays[k]) > 0 {

					newRays = []ray{}

					for i := 0; i < len(rays[k]); i++ {

						newPos, hit, _ = RayMarch(sdf, rays[k][i].pos, rays[k][i].dir)
						numberOfMarches += 1

						if !hit {

							depthStatistics[k] += 1

							if k > 0 {

								norm = CalcNormal(sdf, rays[k][i].pos)

								pixelColor = vector.Sum3(pixelColor, vector.Scale3(colorOfDir(rays[k][i].dir, rays[k][i].pos, j), rays[k][i].weight))

							} else {

								pixelColor = vector.Sum3(pixelColor, colorOfDir(rays[k][i].dir, rays[k][i].pos, j))

							}

						} else {

							norm = CalcNormal(sdf, newPos)

							schlick = Fresnel(rays[k][i].dir, norm, eta1, eta2[j])

							if testReflect && rays[k][i].weight*schlick > FactorCutoff {

								if rays[k][i].inside {
									newRayPos = vector.Sum3(newPos, vector.Scale3(norm, -shift))
								} else {
									newRayPos = vector.Sum3(newPos, vector.Scale3(norm, shift))
								}

								newRays = append(newRays, ray{
									pos:     newRayPos,
									dir:     Reflect(rays[k][i].dir, norm),
									weight:  rays[k][i].weight * schlick,
									rayType: "reflect",
									parent:  i,
									inside:  rays[k][i].inside,
									layer:   k,
								},
								)
							}

							if testRefract {

								refractDir, refract = Refract(rays[k][i].dir, norm, eta1, eta2[j])

								if refract && rays[k][i].weight*(1-schlick) > FactorCutoff {

									if rays[k][i].inside {
										newRayPos = vector.Sum3(newPos, vector.Scale3(norm, shift))
									} else {
										newRayPos = vector.Sum3(newPos, vector.Scale3(norm, -shift))
									}

									newRays = append(newRays, ray{
										pos:     newRayPos,
										dir:     refractDir,
										weight:  rays[k][i].weight * (1 - schlick),
										rayType: "refract",
										parent:  i,
										inside:  !rays[k][i].inside,
										layer:   k,
									})

								}

							}

						}

					}

					k++
					rays = append(rays, newRays)

				}

			}

		}

	}

	pixelColor = vector.Scale3(pixelColor, 1.0/float64(raysPerPixel*raysPerPixel))

	numberOfHits := 3*raysPerPixel*raysPerPixel - depthStatistics[0]

	return pixelColor, numberOfRays, numberOfMarches, numberOfHits, depthStatistics

}

func colorOfDir(dir, pos [3]float64, j int) [3]float64 {

	colors := [6][3]float64{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}

	var newColor [3]float64

	x := 20.0 + pos[0]
	y := x * dir[1] / dir[0]
	z := x * dir[2] / dir[0]

	if math.Abs(y-1.5) < 2.3 && math.Abs(z) < 10.0 {
		newColor = colors[j]
	} else {
		newColor = [3]float64{0, 0, 0}
	}

	return newColor

}
