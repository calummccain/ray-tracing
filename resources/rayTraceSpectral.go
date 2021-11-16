package resources

import (
	"math"

	"github.com/calummccain/coxeter/vector"
)

func RayTraceSpectral(sdf func([3]float64) float64, dir, pos [3]float64, eta1 float64, eta2 []float64, blackBody []float64, iterations int, faces [][]Face, up, left [3]float64, invHeight, invWidth float64, raysPerPixel int, light Light) ([3]float64, int, int, int, []int) {

	spectrum := make([]float64, len(eta2))
	numberOfRays := 0
	numberOfMarches := 0
	depthStatistics := []int{0}

	var wShift, hShift float64
	var shiftedDir [3]float64

	var rays []ray
	var newRays []ray

	var newPos [3]float64
	var newRayPos [3]float64

	var hit bool

	var norm [3]float64
	var schlick float64

	var refractDir [3]float64
	var refract bool

	var j int

	testRefract := true
	testReflect := true

	for l := 0; l < iterations; l++ {
		depthStatistics = append(depthStatistics, 0)
	}

	invRaysPerPixelFloat64 := 1.0 / float64(raysPerPixel)

	for w := 0; w < raysPerPixel; w++ {

		wShift = (float64(w) - float64(raysPerPixel-1)*0.5) * invWidth * invRaysPerPixelFloat64

		for h := 0; h < raysPerPixel; h++ {

			hShift = (float64(h) - float64(raysPerPixel-1)*0.5) * invHeight * invRaysPerPixelFloat64

			shiftedDir = vector.Normalise3(vector.Sum3(vector.Sum3(dir, vector.Scale3(up, hShift)), vector.Scale3(left, wShift)))

			for j = 0; j < len(eta2); j++ {

				numberOfRays += 1

				rays = []ray{
					{pos: pos, dir: shiftedDir, weight: 1, rayType: "initial", parent: 0, inside: false, layer: -1},
				}

				k := 0
				for k < iterations && len(rays) > 0 {

					newRays = []ray{}

					for i := 0; i < len(rays); i++ {

						newPos, hit, _ = RayMarch(sdf, rays[i].pos, rays[i].dir)
						numberOfMarches += 1

						if !hit && !rays[i].inside {

							if rays[i].layer >= -1 {

								norm = CalcNormal(sdf, rays[i].pos)
								spectrum[j] += rays[i].weight * spectralColourLinear(sdf, rays[i].pos, norm, blackBody[j], light)

							}

							continue

						}

						norm = CalcNormal(sdf, newPos)

						schlick = Fresnel(rays[i].dir, norm, eta1/eta2[j])

						if testReflect && rays[i].weight*schlick > FactorCutoff {

							if rays[i].inside {
								newRayPos = vector.Sum3(newPos, vector.Scale3(norm, -RayTraceShift))
							} else {
								newRayPos = vector.Sum3(newPos, vector.Scale3(norm, RayTraceShift))
							}

							newRays = append(newRays, ray{
								pos:     newRayPos,
								dir:     Reflect(rays[i].dir, norm),
								weight:  rays[i].weight * schlick,
								rayType: "reflect",
								parent:  i,
								inside:  rays[i].inside,
								layer:   k,
							},
							)
						}

						if testRefract {

							refractDir, refract = Refract(rays[i].dir, norm, eta1/eta2[j])

							if refract && rays[i].weight*(1-schlick) > FactorCutoff {

								if rays[i].inside {
									newRayPos = vector.Sum3(newPos, vector.Scale3(norm, RayTraceShift))
								} else {
									newRayPos = vector.Sum3(newPos, vector.Scale3(norm, -RayTraceShift))
								}

								newRays = append(newRays, ray{
									pos:     newRayPos,
									dir:     refractDir,
									weight:  rays[i].weight * (1 - schlick),
									rayType: "refract",
									parent:  i,
									inside:  !rays[i].inside,
									layer:   k,
								})

							}

						}

					}

					k++
					rays = newRays

				}

			}

		}

	}

	for j = 0; j < len(eta2); j++ {
		spectrum[j] /= (float64(numberOfRays) * float64(numberOfRays))
	}

	xyz := IntegrateSpectrum(spectrum)
	r, g, b := xyztorgb(SMPTEsystem, xyz[0], xyz[1], xyz[2])
	r, g, b = constrainrgb(r, g, b)
	rgb := normrgb([3]float64{r, g, b})

	// xyz := IntegrateSpectrum(spectrum)
	// xyz = Normalise(xyz)
	// rgb := XYZToRGB(xyz)
	// rgb = [3]float64{GammaCorrection(rgb[0]), GammaCorrection(rgb[1]), GammaCorrection(rgb[2])}

	numberOfHits := raysPerPixel*raysPerPixel - depthStatistics[0]

	return rgb, numberOfRays, numberOfMarches, numberOfHits, depthStatistics

}

func spectralColourRadial(sdf func([3]float64) float64, dir, point, norm [3]float64, blackBody float64) float64 {

	_, hit, t := RayMarch(sdf, point, dir)

	output := 0.0

	if !hit {

		output = blackBody * vector.Dot3(dir, norm) / (t * t)

	}

	return output

}

func spectralColourLinear(sdf func([3]float64) float64, point, norm [3]float64, blackBody float64, light Light) float64 {

	_, hit, t := RayMarch(sdf, point, light.Pos)

	output := 0.0

	if !hit {

		diff := vector.Diff3(point, light.Pos)

		if vector.Dot3(light.Normal, diff)/vector.Dot3(light.Normal, light.Pos) < 0.0 {

			shiftedIntersection := vector.Diff3(diff, vector.Scale3(light.Pos, vector.Dot3(light.Normal, diff)/vector.Dot3(light.Normal, light.Pos)))

			if math.Abs(vector.Dot3(shiftedIntersection, light.Left)) < light.Width && math.Abs(vector.Dot3(shiftedIntersection, light.Up)) < light.Height {

				output = blackBody * math.Abs(vector.Dot3(light.Normal, norm)) / (t * t)

			}

		}

	}

	return output

}
