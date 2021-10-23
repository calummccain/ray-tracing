package resources

import (
	"github.com/calummccain/coxeter/vector"
)

func RayTraceSpectral(sdf func([3]float64) float64, dir, pos [3]float64, eta1 float64, eta2 []float64, blackBody []float64, iterations int, faces [][]Face, up, left [3]float64, invHeight, invWidth float64, raysPerPixel int) ([3]float64, int, int, int, []int) {

	spectrum := make([]float64, 81)
	numberOfRays := 0
	numberOfMarches := 0
	depthStatistics := []int{0}

	var wShift, hShift float64
	var shiftedDir [3]float64

	var rays [][]ray
	var newRays []ray

	var newPos [3]float64
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

	invRaysPerPixelFloat64 := 1.0 / float64(raysPerPixel)

	for w := 0; w < raysPerPixel; w++ {

		wShift = (float64(w) - float64(raysPerPixel-1)*0.5) * invWidth * invRaysPerPixelFloat64

		for h := 0; h < raysPerPixel; h++ {

			hShift = (float64(h) - float64(raysPerPixel-1)*0.5) * invHeight * invRaysPerPixelFloat64

			shiftedDir = vector.Normalise3(vector.Sum3(vector.Sum3(dir, vector.Scale3(up, hShift)), vector.Scale3(left, wShift)))

			for j := 0; j < 81; j += 10 {

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

						newPos, _, _ = RayMarch(sdf, rays[k][i].pos, rays[k][i].dir)
						numberOfMarches += 1

						if !rays[k][i].inside {

							if k > 0 {

								norm = CalcNormal(sdf, rays[k][i].pos)
								spectrum[j] += rays[k][i].weight * spectralColour(sdf, rays[k][i].pos, norm, blackBody[j])

							} else {

								spectrum[j] += spectralColour(sdf, pos, shiftedDir, blackBody[j])

							}

						}

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

					k++
					rays = append(rays, newRays)

				}

			}

		}

	}

	xyz := IntegrateSpectrum(spectrum)
	xyz = Normalise(xyz)
	rgb := XYZToRGB(xyz)
	//rgb = [3]float64{GammaCorrection(rgb[0]), GammaCorrection(rgb[1]), GammaCorrection(rgb[2])}

	numberOfHits := 3*raysPerPixel*raysPerPixel - depthStatistics[0]

	return Desaturate(rgb), numberOfRays, numberOfMarches, numberOfHits, depthStatistics

}

func spectralColour(sdf func([3]float64) float64, point, norm [3]float64, blackBody float64) (output float64) {

	lightPos := [3]float64{-5, 2, 0}
	dir := vector.Normalise3(vector.Diff3(lightPos, point))
	_, hit, _ := RayMarch(sdf, point, dir)

	if !hit {

		output = blackBody * vector.Dot3(dir, norm)

	} else {

		output = 0.0

	}

	return output

}
