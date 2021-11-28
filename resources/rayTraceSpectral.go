package resources

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/calummccain/coxeter/vector"
)

func RayTraceSpectral(sdf func([3]float64) float64, dir, pos [3]float64, eta1 float64, eta2 []float64, blackBody []float64, iterations int, faces [][]Face, up, left [3]float64, invHeight, invWidth float64, raysPerPixel int, lights []Light, sigma float64) (color.RGBA, int, int, int, []int) {

	spectrum := make([]float64, len(eta2))
	hitsPerSpectra := make([]float64, len(eta2))
	numberOfRays := 0
	numberOfMarches := 0
	depthStatistics := []int{0}

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

	var t float64

	var randInt int

	testRefract := true
	testReflect := true

	for l := 0; l < iterations; l++ {
		depthStatistics = append(depthStatistics, 0)
	}

	for w := 0; w < raysPerPixel; w++ {

		shiftedDir = vector.Normalise3(vector.Sum3(vector.Sum3(dir, vector.Scale3(up, invHeight*rand.Float64())), vector.Scale3(left, invWidth*rand.Float64())))

		if raysPerPixel > len(eta2) {
			randInt = w % len(eta2)
		} else {
			randInt = rand.Intn(len(eta2))
		}

		hitsPerSpectra[randInt] += 1.0

		numberOfRays += 1

		rays = []ray{
			{pos: pos, dir: shiftedDir, weight: 1, rayType: "initial", parent: 0, inside: false, layer: -1, length: 0.0},
		}

		k := 0
		for k < iterations && len(rays) > 0 {

			newRays = []ray{}

			for i := 0; i < len(rays); i++ {

				newPos, hit, t = RayMarch(sdf, rays[i].pos, rays[i].dir)
				numberOfMarches += 1

				if !hit && !rays[i].inside {

					if rays[i].layer >= 0 {

						norm = CalcNormal(sdf, rays[i].pos)
						spectrum[randInt] += rays[i].weight * spectralColourLinear(sdf, rays[i].pos, rays[i].dir, lights, rays[i].length, randInt)

					}

					continue

				}

				norm = CalcNormal(sdf, newPos)

				schlick = Fresnel(rays[i].dir, norm, eta1/eta2[randInt])

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
						length:  rays[i].length + t,
					},
					)
				}

				if testRefract {

					refractDir, refract = Refract(rays[i].dir, norm, eta1/eta2[randInt])

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
							length:  rays[i].length + t,
						})

					}

				}

			}

			k++
			rays = newRays

		}

	}

	for i := 0; i < len(eta2); i++ {

		if hitsPerSpectra[i] > 0.0 {

			spectrum[i] /= hitsPerSpectra[i]

		}

	}

	numberOfHits := raysPerPixel*raysPerPixel - depthStatistics[0]

	return SpectrumToRGBA(spectrum, Y_white, sigma), numberOfRays, numberOfMarches, numberOfHits, depthStatistics

}

func spectralColourRadial(sdf func([3]float64) float64, dir, point, norm [3]float64, blackBody float64) float64 {

	_, hit, t := RayMarch(sdf, point, dir)

	output := 0.0

	if !hit {

		output = blackBody * vector.Dot3(dir, norm) / (t * t)

	}

	return output

}

func spectralColourLinear(sdf func([3]float64) float64, point, norm [3]float64, lights []Light, t float64, i int) float64 {

	output := 0.0

	for _, light := range lights {

		diff := vector.Diff3(light.Pos, point)

		if vector.Dot3(light.Normal, diff) >= 0.0 {

			shiftedIntersection := vector.Diff3(vector.Scale3(light.Normal, vector.Dot3(light.Normal, diff)), diff)

			if vector.Dot3(shiftedIntersection, light.Left) < light.Width && vector.Dot3(shiftedIntersection, light.Up) < light.Height &&
				vector.Dot3(shiftedIntersection, light.Left) > -light.Width && vector.Dot3(shiftedIntersection, light.Up) > -light.Height {

				_, hit, tt := RayMarch(sdf, point, light.Normal)

				if (hit && light.Inside) || (!hit && !light.Inside) {

					output += light.Spectrum[i] * math.Abs(vector.Dot3(light.Normal, norm)) / ((t + tt) * (t + tt))

				}

			}

		}

	}

	return output

}
