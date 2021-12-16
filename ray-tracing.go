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
	"sync"
	"time"

	"github.com/calummccain/coxeter/vector"
)

var wg sync.WaitGroup

func main() {

	start := time.Now()

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
	//var r, g, b, a float64

	var camera [3]float64
	var origin [3]float64
	var up [3]float64
	var oc [3]float64
	var left [3]float64

	//var dir [3]float64
	//var colour []color.RGBA

	//colourArray := make([]resources.ColumnColours, configData.ImageConfig.Width)
	// colourArray := make(chan resources.ColumnColours, configData.ImageConfig.Width)

	//var colourRGB [3]float64

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
	//var iFloat float64
	//jFloat float64

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

	if cellData.Metric == 's' || cellData.Metric == 'e' {
		configData.CellGeometryConfig.Model = ""
	}

	sdf := resources.SdfFunction(configData, faceData)

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

	// var numberOfRaysLocal int
	// var numberOfMarchesLocal int
	// var numberOfHitsLocal int
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
		eta2 = append(eta2, 0.5+math.Pow(1.0+float64(i)*invSpectralRaysNumber, 1.2))
		blackBody = append(blackBody, 2000*(5/float64(configData.RaytracingConfig.SpectralRaysNumber))*resources.BlackBodySpectrum(wavelengths[i], configData.LightsConfig[0].Temp))
		whiteLight = append(whiteLight, resources.BlackBodySpectrum(wavelengths[i], 6500))

	}

	lights := []resources.Light{}
	for _, lightConfig := range configData.LightsConfig {
		lights = append(lights, resources.Light{
			Temp:      lightConfig.Temp,
			Intensity: lightConfig.Intensity,
			Pos:       lightConfig.Pos,
			Up:        lightConfig.Up,
			Left:      lightConfig.Left,
			Normal:    lightConfig.Normal,
			Height:    lightConfig.Height,
			Width:     lightConfig.Width,
		})
	}

	for i := 0; i < len(lights); i++ {
		lights[i].SpectrumFromWavelength(wavelengths, configData.RaytracingConfig.SpectralRaysNumber)
		lights[i].LightInside(sdf)
	}

	resources.MergeColourStimulus(configData.RaytracingConfig.SpectralRaysNumber)

	resources.Y_white = resources.IntegrateSpectrum(whiteLight, 1)[1]

	colBlackBody := resources.SpectrumToRGBA(blackBody, resources.Y_white*2000*(5/float64(configData.RaytracingConfig.SpectralRaysNumber)), configData.RaytracingConfig.Sigma)

	spec := make([]float64, len(blackBody))
	red := resources.SpectrumToRGBA(resources.XMatchFunction, 1.0, configData.RaytracingConfig.Sigma)
	green := resources.SpectrumToRGBA(resources.YMatchFunction, 1.0, configData.RaytracingConfig.Sigma)
	blue := resources.SpectrumToRGBA(resources.ZMatchFunction, 1.0, configData.RaytracingConfig.Sigma)

	for k := 0; k < configData.RaytracingConfig.SpectralRaysNumber; k++ {

		spec[k] = 1

		if k > 0 {
			spec[k-1] = 0
		}

		col = resources.SpectrumToRGBA(spec, 1, configData.RaytracingConfig.Sigma)

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

	jobs := make(chan int)

	colourArray := make(chan resources.ColumnColours, configData.ImageConfig.Width)

	for w := 0; w < 4; w++ {

		wg.Add(1)
		go resources.GenerateColumn(
			configData,
			invWidth,
			invHeight,
			jobs,
			oc,
			up,
			left,
			camera,
			sdf,
			faceData,
			eta2,
			blackBody,
			lights,
			colourArray,
			&wg,
		)

	}

	for i := 0; i < configData.ImageConfig.Width; i++ {

		jobs <- i

		//colourArray[i] = resources.ColumnColours{ColumnNumber: i, Colours: colour}

	}

	close(jobs)

	wg.Wait()

	close(colourArray)

	for column := range colourArray {
		for j, colour := range column.Colours {
			rayTracedImage.Set(column.ColumnNumber, j, colour)
		}
	}

	// for _, subArray := range colourArray {
	// 	close(subArray)
	// }
	//close(colourArray)

	//fmt.Print("\033[K\r")

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

	log.Printf("%s", time.Since(start))

}
