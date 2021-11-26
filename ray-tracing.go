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

	"github.com/calummccain/coxeter/data"
	"github.com/calummccain/coxeter/vector"
)

type config struct {
	P               int
	Q               int
	R               float64
	Width           int
	Height          int
	Start           int
	End             int
	Sdf             string
	Cells           []string
	TruncRect       string
	Model           string
	NumberOfFaces   int
	Distance        float64
	Eta1            float64
	Eta2R           float64
	Eta2G           float64
	Eta2B           float64
	NumberOfBounces int
	RaysPerPixel    int
	SphereRadius    float64
	CubeA           float64
	CubeB           float64
	CubeC           float64
	TorusA          float64
	TorusB          float64
	ObjectRotateX   float64
	ObjectRotateY   float64
	ObjectRotateZ   float64
	CameraRotateX   float64
	CameraRotateY   float64
	CameraRotateZ   float64
	Save            bool
	SpectralRays    int
	Temp            int
	Spectral        bool
}

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

	var configData config
	json.Unmarshal(configRead, &configData)

	// Read the heights and width of gif in px
	width := configData.Width
	height := configData.Height

	// Initialise the png images
	spectrumImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{200, 500}})
	rayTracedImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})

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
	cells := configData.Cells

	// Pick and generate the data from the supplier p,q,r and t/r/_
	var cellData data.CellData
	cell := [2]int{configData.P, configData.Q}
	if configData.TruncRect == "r" {
		switch cell {
		case [2]int{3, 3}:
			cellData = data.RectifiedTetrahedronData(configData.R)
		case [2]int{3, 4}:
			cellData = data.RectifiedOctahedronData(configData.R)
		case [2]int{3, 5}:
			cellData = data.RectifiedIcosahedronData(configData.R)
		case [2]int{4, 3}:
			cellData = data.RectifiedHexahedronData(configData.R)
		case [2]int{5, 3}:
			cellData = data.RectifiedDodecahedronData(configData.R)
		default:
			cellData = data.TetrahedronData(configData.R)
		}
	} else if configData.TruncRect == "t" {
		switch cell {
		case [2]int{3, 3}:
			cellData = data.TruncatedTetrahedronData(configData.R)
		case [2]int{3, 4}:
			cellData = data.TruncatedOctahedronData(configData.R)
		case [2]int{3, 5}:
			cellData = data.TruncatedIcosahedronData(configData.R)
		case [2]int{4, 3}:
			cellData = data.TruncatedHexahedronData(configData.R)
		case [2]int{5, 3}:
			cellData = data.TruncatedDodecahedronData(configData.R)
		default:
			cellData = data.TetrahedronData(configData.R)
		}
	} else {
		switch cell {
		case [2]int{3, 3}:
			cellData = data.TetrahedronData(configData.R)
		case [2]int{3, 4}:
			cellData = data.OctahedronData(configData.R)
		case [2]int{3, 5}:
			cellData = data.IcosahedronData(configData.R)
		case [2]int{3, 6}:
			cellData = data.TriangularData(configData.R, configData.NumberOfFaces)
			cellData.C = [4]float64{1, -4, 0, 0}
		case [2]int{4, 3}:
			cellData = data.HexahedronData(configData.R)
		case [2]int{4, 4}:
			cellData = data.SquareData(configData.R, configData.NumberOfFaces)
		case [2]int{5, 3}:
			cellData = data.DodecahedronData(configData.R)
		case [2]int{6, 3}:
			cellData = data.HexagonalData(configData.R, configData.NumberOfFaces)
		default:
			cellData = data.HyperbolicData(configData.P, configData.Q, configData.R, configData.NumberOfFaces)
		}
	}

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
			faceData = append(faceData, resources.GenerateFacesHyperbolic(cellData.NumFaces, cellData.Faces, vertexData, cellData.Metric, cellData.Vv, configData.Model, cellData.Matrices.F(vector.TransformVertices([][4]float64{cellData.C}, cells[i], cellData.Matrices)[0])))

		}

	}

	for i := 0; i < len(faceData); i++ {

		for j := 0; j < len(faceData[i]); j++ {

			faceData[i][j].SphereCenter = resources.RotateXYZ(faceData[i][j].SphereCenter, configData.ObjectRotateX, configData.ObjectRotateY, configData.ObjectRotateZ)
			faceData[i][j].Normal = resources.RotateXYZ(faceData[i][j].Normal, configData.ObjectRotateX, configData.ObjectRotateY, configData.ObjectRotateZ)

		}

	}

	type sdfFunction func([3]float64) float64

	var sdf sdfFunction

	if configData.Sdf == "sphere" {
		sdf = func(p [3]float64) float64 {
			return resources.SdfSphere(resources.RotateXYZ(p, configData.ObjectRotateX, configData.ObjectRotateY, configData.ObjectRotateZ), configData.SphereRadius)
		}
	} else if configData.Sdf == "cube" {
		sdf = func(p [3]float64) float64 {
			return resources.SdfCube(resources.RotateXYZ(p, configData.ObjectRotateX, configData.ObjectRotateY, configData.ObjectRotateZ), configData.CubeA, configData.CubeB, configData.CubeC)
		}
	} else if configData.Sdf == "torus" {
		sdf = func(p [3]float64) float64 {
			return resources.SdfTorus(resources.RotateXYZ(p, configData.ObjectRotateX, configData.ObjectRotateY, configData.ObjectRotateZ), configData.TorusA, configData.TorusB)
		}
	} else if configData.Sdf == "spheres" {
		sdf = func(p [3]float64) float64 {
			return resources.SdfSpheres(resources.RotateXYZ(p, configData.ObjectRotateX, configData.ObjectRotateY, configData.ObjectRotateZ))
		}
	} else if configData.Sdf == "seh" {

		if cellData.Metric == 's' || cellData.Metric == 'e' {
			configData.Model = ""
		}

		sdf = func(p [3]float64) float64 {
			return resources.Sdf(p, faceData, configData.Model)
		}

	}

	light := resources.Light{Pos: [3]float64{0, -10, 0}, Up: [3]float64{1, 0, 0}, Left: [3]float64{0, 0, 1}, Normal: [3]float64{0, -1, 0}, Height: 1.0, Width: 1.0}
	//light.Normal = vector.Cross3(light.Up, light.Left)

	camera = resources.RotateXYZ([3]float64{configData.Distance, 0, 0}, configData.CameraRotateX, configData.CameraRotateY, configData.CameraRotateZ)
	origin = [3]float64{0, 0, 0}
	up = resources.RotateXYZ([3]float64{0, 0, 1}, configData.CameraRotateX, configData.CameraRotateY, configData.CameraRotateZ)

	oc = vector.Diff3(origin, camera)
	oc = vector.Normalise3(oc)
	left = vector.Cross3(oc, up)

	invWidth := 1.0 / float64(width)
	invHeight := 1.0 / float64(height)

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

	for i := 0; i < configData.SpectralRays; i++ {

		wavelengths = append(wavelengths, 380.0+400*(float64(i)+0.5)/float64(configData.SpectralRays))
		eta2 = append(eta2, 1.0+math.Pow(1.0+float64(i)/float64(configData.SpectralRays), 1.2))
		blackBody = append(blackBody, 10000*resources.BlackBodySpectrum(wavelengths[i], float64(configData.Temp)))
		whiteLight = append(whiteLight, resources.BlackBodySpectrum(wavelengths[i], 6500))

	}

	newX := make([]float64, configData.SpectralRays)
	newY := make([]float64, configData.SpectralRays)
	newZ := make([]float64, configData.SpectralRays)

	for ii := 0; ii < configData.SpectralRays; ii++ {

		for jj := 0; jj < 80/configData.SpectralRays; jj++ {

			//fmt.Println(ii*80/configData.SpectralRays + jj)

			newX[ii] += resources.XMatchFunction[ii*80/configData.SpectralRays+jj]
			newY[ii] += resources.YMatchFunction[ii*80/configData.SpectralRays+jj]
			newZ[ii] += resources.ZMatchFunction[ii*80/configData.SpectralRays+jj]

		}

		newX[ii] *= float64(configData.SpectralRays) / 80.0
		newY[ii] *= float64(configData.SpectralRays) / 80.0
		newZ[ii] *= float64(configData.SpectralRays) / 80.0

	}

	// fmt.Println(newX)
	// fmt.Println(newY)
	// fmt.Println(newZ)

	resources.XMatchFunction = newX
	resources.YMatchFunction = newY
	resources.ZMatchFunction = newZ

	resources.Y_white = resources.IntegrateSpectrum(whiteLight, 1)[1]

	// rr, gg, bb := resources.Constrainrgb(
	// 	resources.IntegrateSpectrum(blackBody, resources.Y_white)[0],
	// 	resources.IntegrateSpectrum(blackBody, resources.Y_white)[1],
	// 	resources.IntegrateSpectrum(blackBody, resources.Y_white)[2],
	// )

	// fmt.Println(resources.Y_white)
	// fmt.Println(blackBody)
	// fmt.Println(resources.IntegrateSpectrum(blackBody, resources.Y_white))
	// fmt.Println(rr, gg, bb)
	// fmt.Println(resources.XYZToRGB(rr, gg, bb))

	colBlackBody := resources.SpectrumToRGBA(blackBody, resources.Y_white*10000)

	//fmt.Println(resources.Y_white)

	spec := make([]float64, len(blackBody))
	red := resources.SpectrumToRGBA(resources.XMatchFunction, 1.0)

	green := resources.SpectrumToRGBA(resources.YMatchFunction, 1.0)
	// fmt.Println(resources.IntegrateSpectrum(resources.YMatchFunction, 1.0))
	// fmt.Println(resources.Constrainrgb(
	// 	resources.IntegrateSpectrum(resources.YMatchFunction, 1.0)[0],
	// 	resources.IntegrateSpectrum(resources.YMatchFunction, 1.0)[1],
	// 	resources.IntegrateSpectrum(resources.YMatchFunction, 1.0)[2],
	// ))
	// fmt.Println(green)
	blue := resources.SpectrumToRGBA(resources.ZMatchFunction, 1.0)

	for k := 0; k < configData.SpectralRays; k++ {

		spec[k] = 1

		if k > 0 {
			spec[k-1] = 0
		}

		col = resources.SpectrumToRGBA(spec, 1)

		for i := (200 * k) / configData.SpectralRays; i < (200*(k+1))/configData.SpectralRays; i++ {
			for j := 0; j < 100; j++ {
				spectrumImage.Set(i, j, col)
			}
		}

		for i := (200 * k) / configData.SpectralRays; i < (200*(k+1))/configData.SpectralRays; i++ {
			for j := 0; j < int(50*resources.XMatchFunction[k]); j++ {
				spectrumImage.Set(i, 100+j, red)
			}
		}

		for i := (200 * k) / configData.SpectralRays; i < (200*(k+1))/configData.SpectralRays; i++ {
			for j := 0; j < int(50*resources.YMatchFunction[k]); j++ {
				spectrumImage.Set(i, 200+j, green)
			}
		}

		for i := (200 * k) / configData.SpectralRays; i < (200*(k+1))/configData.SpectralRays; i++ {
			for j := 0; j < int(50*resources.ZMatchFunction[k]); j++ {
				spectrumImage.Set(i, 300+j, blue)
			}
		}

		for i := (200 * k) / configData.SpectralRays; i < (200*(k+1))/configData.SpectralRays; i++ {
			for j := 0; j < int(50*blackBody[k]/1e14); j++ {
				spectrumImage.Set(i, 400+j, colBlackBody)
			}
		}

	}

	f, _ := os.Create("images/spectrum.png")
	png.Encode(f, spectrumImage)

	for i := 0; i < width; i++ {

		iFloat = float64(i)*invWidth - 0.5*(1-invWidth)

		for j := 0; j < height; j++ {

			fmt.Print("\033[K\r")
			fmt.Print(i, " - ", j)

			jFloat = float64(j)*invHeight - 0.5*(1-invHeight)

			dir = vector.Sum3(vector.Sum3(oc, vector.Scale3(up, jFloat)), vector.Scale3(left, iFloat))

			if !configData.Spectral {

				colourRGB, numberOfRaysLocal, numberOfMarchesLocal, numberOfHitsLocal, _ = resources.RayTrace(
					sdf,
					dir,
					camera,
					configData.Eta1,
					[3]float64{configData.Eta2R, configData.Eta2G, configData.Eta2B},
					configData.NumberOfBounces,
					faceData,
					up,
					left,
					invHeight,
					invWidth,
					configData.RaysPerPixel,
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
					configData.Eta1,
					eta2,
					blackBody,
					configData.NumberOfBounces,
					faceData,
					up,
					left,
					invHeight,
					invWidth,
					configData.RaysPerPixel,
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

	if configData.Save {

		f, _ := os.Create(fmt.Sprintf("images/png/%d.png", timeString))
		png.Encode(f, rayTracedImage)

	}

	if configData.Save {

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
