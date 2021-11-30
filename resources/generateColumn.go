package resources

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/calummccain/coxeter/vector"
)

func GenerateColumn(configData Config, invWidth, invHeight float64, jobs <-chan int, oc, up, left, camera [3]float64, sdf func([3]float64) float64, faceData [][]Face, eta2 []float64, blackbody []float64, lights []Light, colourChannel chan<- ColumnColours, wg *sync.WaitGroup) {

	for job := range jobs {

		fmt.Println(job)

		var jFloat float64
		iFloat := float64(job)*invWidth - 0.5*(1-invWidth)

		var dir [3]float64
		var colourRGB [3]float64
		//var numberOfRaysLocal, numberOfMarchesLocal, numberOfHitsLocal int
		var r, g, b, a float64
		var colour color.RGBA

		column := make([]color.RGBA, configData.ImageConfig.Height)

		for j := 0; j < configData.ImageConfig.Height; j++ {

			// fmt.Print("\033[K\r")
			// fmt.Print(job, " - ", j)

			jFloat = float64(j)*invHeight - 0.5*(1-invHeight)

			dir = vector.Sum3(vector.Sum3(oc, vector.Scale3(up, jFloat)), vector.Scale3(left, iFloat))

			if !configData.RaytracingConfig.Spectral {

				colourRGB, _, _, _, _ = RayTrace(
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

				colour, _, _, _, _ = RayTraceSpectral(
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

		colourChannel <- ColumnColours{ColumnNumber: job, Colours: column}
		//fmt.Println(<-colourChannel)

	}

	wg.Done()

	//return column, numberOfRaysLocal, numberOfMarchesLocal, numberOfHitsLocal

}
