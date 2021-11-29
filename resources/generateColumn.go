package resources

import (
	"fmt"
	"image/color"

	"github.com/calummccain/coxeter/vector"
)

func GenerateColumn(configData Config, invWidth, invHeight float64, i int, oc, up, left, camera [3]float64, iFloat float64, sdf func([3]float64) float64, faceData [][]Face, eta2 []float64, blackbody []float64, lights []Light) ([]color.Color, int, int, int) {

	var jFloat float64
	var dir [3]float64
	var colourRGB [3]float64
	var numberOfRaysLocal, numberOfMarchesLocal, numberOfHitsLocal int
	var r, g, b, a float64
	var colour color.Color

	column := make([]color.Color, configData.ImageConfig.Height)

	for j := 0; j < configData.ImageConfig.Height; j++ {

		fmt.Print("\033[K\r")
		fmt.Print(i, " - ", j)

		jFloat = float64(j)*invHeight - 0.5*(1-invHeight)

		dir = vector.Sum3(vector.Sum3(oc, vector.Scale3(up, jFloat)), vector.Scale3(left, iFloat))

		if !configData.RaytracingConfig.Spectral {

			colourRGB, numberOfRaysLocal, numberOfMarchesLocal, numberOfHitsLocal, _ = RayTrace(
				sdf,
				dir,
				camera,
				configData.MaterialConfig.Eta1,
				[3]float64{configData.MaterialConfig.Eta2R, configData.MaterialConfig.Eta2G, configData.MaterialConfig.Eta2B},
				configData.RaytracingConfig.NumberOfBounces,
				faceData,
				up,
				left,
				invHeight,
				invWidth,
				configData.RaytracingConfig.RaysPerPixel,
			)

			r = colourRGB[0]
			g = colourRGB[1]
			b = colourRGB[2]
			a = 255

			r = 255 * r
			g = 255 * g
			b = 255 * b

			colour = color.RGBA{
				uint8(r),
				uint8(g),
				uint8(b),
				uint8(a),
			}

		} else {

			colour, numberOfRaysLocal, numberOfMarchesLocal, numberOfHitsLocal, _ = RayTraceSpectral(
				sdf,
				dir,
				camera,
				configData.MaterialConfig.Eta1,
				eta2,
				blackbody,
				configData.RaytracingConfig.NumberOfBounces,
				faceData,
				up,
				left,
				invHeight,
				invWidth,
				configData.RaytracingConfig.RaysPerPixel,
				lights,
				configData.RaytracingConfig.Sigma,
			)

		}

		column[j] = colour

		//depthStatistics = resources.SumInt(depthStatistics, depthStatisticsLocal)

	}

	return column, numberOfRaysLocal, numberOfMarchesLocal, numberOfHitsLocal

}
