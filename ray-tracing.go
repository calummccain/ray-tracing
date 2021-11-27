package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"ray-tracing/resources"
	"runtime/pprof"
	"time"

	"github.com/calummccain/coxeter/vector"
)

func main() {

	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	timeString := time.Now().Unix()

	configJson, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}

	defer configJson.Close()

	configRead, _ := ioutil.ReadAll(configJson)

	var configData resources.Config
	json.Unmarshal(configRead, &configData)

	// Initialise the png images
	spectrumImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{200, 500}})
	rayTracedImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{configData.ImageConfig.Width, configData.ImageConfig.Height}})

	var col color.Color
	var r, g, b, a float64

	var camera [3]float64
	var origin [3]float64
	var up [3]float64
	var oc [3]float64
	var left [3]float64

	var dir [3]float64
	var colour color.RGBA

	var colourRGB [3]float64

	// Local coordinates for view frame
	//               1
	//               |
	//         jFloat|--------*
	//               |        |
	// -1------------+--------+--- 1
	//               |        iFloat
	//               |
	//               |
	//               -1
	var iFloat, jFloat float64

	// Cells is a list of the names of the cells to calculate for
	cells := configData.CellGeometryConfig.Cells

	// Pick and generate the data from the supplier p,q,r and t/r/_
	cellData := resources.SelectGeometry(configData.CellGeometryConfig)

	faceData := [][]resources.Face{}

	if cellData.Metric == 's' {

		var vertexData []resources.VertexSpherical

		for i := 0; i < len(cells); i++ {

			vertexData = resources.GenerateVerticesSpherical(cellData.Vertices, cells[i], cellData.Matrices, cellData.NumVertices)
			faceData = append(faceData, resources.GenerateFacesSpherical(cellData.NumFaces, cellData.Faces, vertexData, cellData.Matrices.F(vector.TransformVertices([][4]float64{cellData.C}, cells[i], cellData.Matrices)[0])))

		}

	} else if cellData.Metric == 'e' {

		var vertexData []resources.VertexEuclidean

		for i := 0; i < len(cells); i++ {

			vertexData = resources.GenerateVerticesEuclidean(cellData.Vertices, cells[i], cellData.Matrices, cellData.NumVertices)
			faceData = append(faceData, resources.GenerateFacesEuclidean(cellData.NumFaces, cellData.Faces, vertexData, cellData.Matrices.F(vector.TransformVertices([][4]float64{cellData.C}, cells[i], cellData.Matrices)[0])))

		}

	} else {

		var vertexData []resources.VertexHyperbolic

		for i := 0; i < len(cells); i++ {

			vertexData = resources.GenerateVerticesHyperbolic(cellData.Vertices, cells[i], cellData.Matrices, cellData.NumVertices)
			faceData = append(faceData, resources.GenerateFacesHyperbolic(cellData.NumFaces, cellData.Faces, vertexData, cellData.Metric, cellData.Vv, configData.CellGeometryConfig.Model, cellData.Matrices.F(vector.TransformVertices([][4]float64{cellData.C}, cells[i], cellData.Matrices)[0])))

		}

	}

	for i := 0; i < len(faceData); i++ {

		for j := 0; j < len(faceData[i]); j++ {

			faceData[i][j].SphereCenter = resources.RotateXYZ(
				faceData[i][j].SphereCenter, configData.ObjectConfig.ObjectRotateX,
				configData.ObjectConfig.ObjectRotateY,
				configData.ObjectConfig.ObjectRotateZ,
			)
			faceData[i][j].Normal = resources.RotateXYZ(
				faceData[i][j].Normal,
				configData.ObjectConfig.ObjectRotateX,
				configData.ObjectConfig.ObjectRotateY,
				configData.ObjectConfig.ObjectRotateZ,
			)

		}

	}

	sdf := resources.SdfFunction(configData, faceData)

	light := resources.Light{Pos: [3]float64{0, -10, 0}, Up: [3]float64{1, 0, 0}, Left: [3]float64{0, 0, 1}, Normal: [3]float64{0, -1, 0}, Height: 0.4, Width: 0.4}

	camera = resources.RotateXYZ(
		[3]float64{configData.CameraConfig.Distance, 0, 0},
		configData.CameraConfig.CameraRotateX,
		configData.CameraConfig.CameraRotateY,
		configData.CameraConfig.CameraRotateZ,
	)
	origin = [3]float64{0, 0, 0}
	up = resources.RotateXYZ(
		[3]float64{0, 0, 1},
		configData.CameraConfig.CameraRotateX,
		configData.CameraConfig.CameraRotateY,
		configData.CameraConfig.CameraRotateZ,
	)

	oc = vector.Diff3(origin, camera)
	oc = vector.Normalise3(oc)
	left = vector.Cross3(oc, up)

	invWidth := 1.0 / float64(configData.ImageConfig.Width)
	invHeight := 1.0 / float64(configData.ImageConfig.Height)

	numberOfRays := 0
	numberOfMarches := 0
	averageDepth := 0.0
	hitPixels := 0
	//depthStatistics := []int{}

	var numberOfRaysLocal int
	var numberOfMarchesLocal int
	var numberOfHitsLocal int
	//var depthStatisticsLocal []int

	//depthStatistics := make([]int, configData.NumberOfBounces+1)

	// for l := 0; l <= configData.NumberOfBounces; l++ {
	// 	depthStatistics = append(depthStatistics, 0)
	// }

	wavelengths := []float64{}
	eta2 := []float64{}
	blackBody := []float64{}
	whiteLight := []float64{}

	invSpectralRaysNumber := 1.0 / float64(configData.RaytracingConfig.SpectralRaysNumber)

	for i := 0; i < configData.RaytracingConfig.SpectralRaysNumber; i++ {

		wavelengths = append(wavelengths, 380.0+400*(float64(i)+0.5)*invSpectralRaysNumber)
		eta2 = append(eta2, 1.0+math.Pow(1.0+float64(i)*invSpectralRaysNumber, 1.2))
		blackBody = append(blackBody, 10000*resources.BlackBodySpectrum(wavelengths[i], configData.LightsConfig[0].Temp))
		whiteLight = append(whiteLight, resources.BlackBodySpectrum(wavelengths[i], 6500))

	}

	resources.MergeColourStimulus(configData.RaytracingConfig.SpectralRaysNumber)

	resources.Y_white = resources.IntegrateSpectrum(whiteLight, 1)[1]

	colBlackBody := resources.SpectrumToRGBA(blackBody, resources.Y_white*10000)

	spec := make([]float64, len(blackBody))
	red := resources.SpectrumToRGBA(resources.XMatchFunction, 1.0)
	green := resources.SpectrumToRGBA(resources.YMatchFunction, 1.0)
	blue := resources.SpectrumToRGBA(resources.ZMatchFunction, 1.0)

	for k := 0; k < configData.RaytracingConfig.SpectralRaysNumber; k++ {

		spec[k] = 1

		if k > 0 {
			spec[k-1] = 0
		}

		col = resources.SpectrumToRGBA(spec, 1)

		for i := (200 * k) / configData.RaytracingConfig.SpectralRaysNumber; i < (200*(k+1))/configData.RaytracingConfig.SpectralRaysNumber; i++ {

			for j := 0; j < 100; j++ {
				spectrumImage.Set(i, j, col)
			}

			for j := 0; j < int(50*resources.XMatchFunction[k]); j++ {
				spectrumImage.Set(i, 100+j, red)
			}

			for j := 0; j < int(50*resources.YMatchFunction[k]); j++ {
				spectrumImage.Set(i, 200+j, green)
			}

			for j := 0; j < int(50*resources.ZMatchFunction[k]); j++ {
				spectrumImage.Set(i, 300+j, blue)
			}

			for j := 0; j < int(50*blackBody[k]/1e14); j++ {
				spectrumImage.Set(i, 400+j, colBlackBody)
			}

		}

	}

	f, _ := os.Create("images/spectrum.png")
	png.Encode(f, spectrumImage)

	for i := 0; i < configData.ImageConfig.Width; i++ {

		iFloat = float64(i)*invWidth - 0.5*(1-invWidth)

		for j := 0; j < configData.ImageConfig.Height; j++ {

			fmt.Print("\033[K\r")
			fmt.Print(i, " - ", j)

			jFloat = float64(j)*invHeight - 0.5*(1-invHeight)

			dir = vector.Sum3(vector.Sum3(oc, vector.Scale3(up, jFloat)), vector.Scale3(left, iFloat))

			if !configData.RaytracingConfig.Spectral {

				colourRGB, numberOfRaysLocal, numberOfMarchesLocal, numberOfHitsLocal, _ = resources.RayTrace(
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

				colour, numberOfRaysLocal, numberOfMarchesLocal, numberOfHitsLocal, _ = resources.RayTraceSpectral(
					sdf,
					dir,
					camera,
					configData.MaterialConfig.Eta1,
					eta2,
					blackBody,
					configData.RaytracingConfig.NumberOfBounces,
					faceData,
					up,
					left,
					invHeight,
					invWidth,
					configData.RaytracingConfig.RaysPerPixel,
					light,
				)

			}

			rayTracedImage.Set(i, j, colour)

			numberOfRays += numberOfRaysLocal
			numberOfMarches += numberOfMarchesLocal
			hitPixels += numberOfHitsLocal
			//depthStatistics = resources.SumInt(depthStatistics, depthStatisticsLocal)

		}

	}

	fmt.Print("\033[K\r")

	averageDepth = float64(numberOfMarches-numberOfRays+hitPixels) / float64(hitPixels)

	fmt.Println("Number of Rays: ", numberOfRays)
	fmt.Println("Number of Marches: ", numberOfMarches)
	fmt.Println("Number of hit Pixels: ", hitPixels)
	fmt.Println("Average Depth of Ray : ", math.Log2(averageDepth))
	//fmt.Println("Depth Statistics: ", depthStatistics)
	fmt.Println("Number of Faces: ", len(faceData[0]))

	f, _ = os.Create("images/test.png")
	png.Encode(f, rayTracedImage)

	if configData.ImageConfig.Save {

		f, _ := os.Create(fmt.Sprintf("images/png/%d.png", timeString))
		png.Encode(f, rayTracedImage)

		origJson, err := os.Open("config.json")
		if err != nil {
			fmt.Println(err)
		}

		copyJson, err := os.Create(fmt.Sprintf("images/data/%d.json", timeString))
		if err != nil {
			fmt.Println(err)
		}

		_, err = io.Copy(copyJson, origJson)
		if err != nil {
			fmt.Println(err)
		}

		copyJson.Close()

		origJson.Close()

	}

}
