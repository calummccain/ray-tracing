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

func RayTrace(dir, pos [3]float64, eta1 float64, eta2 [6]float64, iterations int, faces [][]FaceHyperbolic, up, left [3]float64, invHeight, invWidth float64, raysPerPixel int) [3]float64 {

	pixelColor := [3]float64{0, 0, 0}

	var shiftedDir [3]float64

	var wShift, hShift float64

	var newPos [3]float64
	var newRayPos [3]float64
	var hit bool

	var norm [3]float64

	var newRays []ray

	var schlick float64

	var refractDir [3]float64
	var refract bool

	testRefract := true
	testReflect := true

	shift := RayTraceShift

	for w := 0; w < raysPerPixel; w++ {

		wShift = (float64(w) - float64(raysPerPixel-1)*0.5) * invWidth / float64(raysPerPixel)

		for h := 0; h < raysPerPixel; h++ {

			hShift = (float64(h) - float64(raysPerPixel-1)*0.5) * invHeight / float64(raysPerPixel)

			shiftedDir = vector.Normalise3(vector.Sum3(vector.Sum3(dir, vector.Scale3(up, hShift)), vector.Scale3(left, wShift)))

			for j := 0; j < 6; j += 2 {

				// fmt.Println("///////////////////////////")

				rays := [][]ray{
					{
						{pos: pos, dir: shiftedDir, weight: 1, rayType: "initial", parent: 0, inside: false, layer: -1},
					},
				}

				k := 0
				for k < iterations && len(rays[k]) > 0 {

					newRays = []ray{}

					for i := 0; i < len(rays[k]); i++ {

						newPos, hit, _ = RayMarch(rays[k][i].pos, rays[k][i].dir, faces)

						if !hit {

							if k > 0 {

								// if rays[k][i].inside {

								// 	fmt.Println(rays[k][i])

								// }

								norm = CalcNormal(faces, rays[k][i].pos)

								pixelColor = vector.Sum3(pixelColor, vector.Scale3(colorOfDir(rays[k][i].dir, rays[k][i].pos, j), rays[k][i].weight))

							} else {

								//pixelColor[j] += colorOfDir2(pos, shiftedDir)[j]
								pixelColor = vector.Sum3(pixelColor, colorOfDir(rays[k][i].dir, rays[k][i].pos, j))

							}

						} else {

							norm = CalcNormal(faces, newPos)

							// fmt.Println("inside: ", rays[k][i].inside)

							schlick = Fresnel(rays[k][i].dir, norm, eta1, eta2[j])

							//fmt.Println("schlick: ", schlick)

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

	return pixelColor

}

func colorOfDir(dir, pos [3]float64, j int) [3]float64 {

	// a := math.Abs(dir[0])
	// b := math.Abs(dir[1])
	// c := math.Abs(dir[2])

	colors := [6][3]float64{
		{1, 0, 0},
		{1, 1, 0},
		{0, 1, 0},
		{0, 1, 1},
		{0, 0, 1},
		{1, 0, 1},
	}

	var newColor [3]float64
	// if a >= b && a >= c && dir[0] >= 0 {

	// 	newColor = [3]float64{0, 0, 1}

	// } else if a >= b && a >= c && dir[0] <= 0 && dir[1] >= 0 && dir[2] >= 0 {

	// 	newColor = [3]float64{1, 0.5, 0}

	// } else if a >= b && a >= c && dir[0] <= 0 && dir[1] <= 0 && dir[2] <= 0 {

	// 	newColor = [3]float64{0.5, 1, 0}

	// } else if a >= b && a >= c && dir[0] <= 0 && dir[1] <= 0 && dir[2] >= 0 {

	// 	newColor = [3]float64{0.5, 0.5, 0}

	// } else if a >= b && a >= c && dir[0] <= 0 && dir[1] >= 0 && dir[2] <= 0 {

	// 	newColor = [3]float64{1, 1, 0}

	// } else if b >= a && b >= c && dir[1] >= 0 {

	// 	newColor = [3]float64{1, 0, 0}

	// } else if b >= a && b >= c && dir[1] <= 0 {

	// 	newColor = [3]float64{0, 1, 1}

	// } else if c >= a && c >= b && dir[2] >= 0 {

	// 	newColor = [3]float64{0, 1, 0}

	// } else {

	// 	newColor = [3]float64{1, 0, 1}

	// }

	// if dir[1] < 0.025 && dir[1] > -0.025 {
	// 	newColor = colors[j]
	// } else if dir[1] < 0.075 && dir[1] > -0.075 {
	// 	newColor = [3]float64{0, 0, 0}
	// } else if dir[1] < 0.125 && dir[1] > -0.125 {
	// 	newColor = colors[j]
	// } else {
	// 	newColor = [3]float64{0, 0, 0}
	// }
	x := 10.0 + pos[0]
	y := x * dir[1] / dir[0]
	z := x * dir[2] / dir[0]

	if math.Abs(y-1.5) < 2.3 && math.Abs(z) < 10.0 {
		newColor = colors[j]
	} else {
		newColor = [3]float64{0, 0, 0}
	}

	// if dir[1] < 0.015 && dir[1] > -0.015 {
	// 	newColor = colors[j]
	// } else {
	// 	newColor = [3]float64{0, 0, 0}
	// }

	// n := 0.5 * (dir[1] + 1)
	// newColor = [3]float64{n, n, n}

	return newColor

}

// func colorOfDir(pos, norm [3]float64, faces [][]FaceHyperbolic) [3]float64 {

// 	I := 100.0

// 	dir := vector.Normalise3(vector.Diff3([3]float64{-2, 0, 1}, pos))

// 	_, hit, t := RayMarch(pos, dir, faces)

// 	col := [3]float64{0, 0, 0}

// 	var dot float64

// 	if !hit {

// 		dot = math.Min(I*math.Abs(vector.Dot3(dir, norm))/(t*t), 1)

// 		col = vector.Sum3(col, [3]float64{dot, dot, dot})

// 	} else {

// 		col = vector.Sum3(col, [3]float64{0, 0, 0})

// 	}

// 	dir = vector.Normalise3(vector.Diff3([3]float64{-2, 0.866, -0.5}, pos))

// 	_, hit, t = RayMarch(pos, dir, faces)

// 	if !hit {

// 		dot = math.Min(I*math.Abs(vector.Dot3(dir, norm))/(t*t), 1)

// 		col = vector.Sum3(col, [3]float64{dot, dot, dot})

// 	} else {

// 		col = vector.Sum3(col, [3]float64{0, 0, 0})

// 	}

// 	dir = vector.Normalise3(vector.Diff3([3]float64{-2, -0.866, -0.5}, pos))

// 	_, hit, t = RayMarch(pos, dir, faces)

// 	if !hit {

// 		dot = math.Min(I*math.Abs(vector.Dot3(dir, norm))/(t*t), 1)

// 		col = vector.Sum3(col, [3]float64{dot, dot, dot})

// 	} else {

// 		col = vector.Sum3(col, [3]float64{0, 0, 0})

// 	}

// 	return col

// }

// func colorOfDir2(pos, norm [3]float64) [3]float64 {

// 	dir := vector.Normalise3(vector.Diff3([3]float64{-2, 0, 1}, pos))

// 	col := [3]float64{0, 0, 0}

// 	var dot float64

// 	dot = math.Abs(vector.Dot3(dir, norm))

// 	col = vector.Sum3(col, [3]float64{0, dot, 0})

// 	dir = vector.Normalise3(vector.Diff3([3]float64{-2, 0.866, -0.5}, pos))

// 	dot = math.Abs(vector.Dot3(dir, norm))

// 	col = vector.Sum3(col, [3]float64{0, 0, dot})

// 	dir = vector.Normalise3(vector.Diff3([3]float64{-2, -0.866, -0.5}, pos))

// 	dot = math.Abs(vector.Dot3(dir, norm))

// 	col = vector.Sum3(col, [3]float64{dot, 0, 0})

// 	return col

// }
