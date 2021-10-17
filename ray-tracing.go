package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"math"
	"os"
	"ray-tracing/resources"
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
	Cells           []string
	TruncRect       string
	Model           string
	NumberOfFaces   int
	Distance        float64
	Angle           float64
	Eta1            float64
	Eta2R           float64
	Eta2G           float64
	Eta2B           float64
	NumberOfBounces int
	RaysPerPixel    int
}

func main() {

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

	// Start and end determine number of frames
	start := configData.Start
	end := configData.End

	// Initialise the png image
	upleft := image.Point{0, 0}
	downRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upleft, downRight})

	var col color.Color
	var r, g, b, a float64
	// Variables for the camera
	//
	//
	//
	//
	//
	//
	//
	//
	//
	var camera [3]float64
	var origin [3]float64
	var up [3]float64
	var oc [3]float64
	var left [3]float64

	// Variables for animation
	//  - angle: angle of rotation
	var angle float64
	var dir [3]float64
	var colour [3]float64

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

	var vertexData []resources.VertexHyperbolic
	faceData := [][]resources.FaceHyperbolic{}

	for i := 0; i < len(cells); i++ {

		vertexData = resources.GenerateVerticesHyperbolic(cellData.Vertices, cells[i], cellData.Matrices, cellData.NumVertices)

		faceData = append(faceData, resources.GenerateFacesHyperbolic(cellData.NumFaces, cellData.Faces, vertexData, cellData.Metric, cellData.Vv, "poincare", cellData.Matrices.F(vector.TransformVertices([][4]float64{cellData.C}, cells[i], cellData.Matrices)[0])))

	}

	for time := start; time < end; time++ {

		angle = configData.Angle

		camera = [3]float64{configData.Distance * math.Cos(angle),
			configData.Distance * math.Sin(angle) * data.Rt_2,
			configData.Distance * math.Sin(angle) * data.Rt_2,
		}
		origin = [3]float64{0, 0, 0}
		up = [3]float64{
			0,
			-math.Sin(angle),
			math.Cos(angle),
		}

		// camera = [3]float64{configData.Distance * math.Cos(angle),
		// 	configData.Distance * math.Sin(angle),
		// 	0,
		// }
		// origin = [3]float64{0, 0, 0}
		// up = [3]float64{
		// 	0,
		// 	0,
		// 	1,
		// }

		oc = vector.Normalise3(vector.Diff3(origin, camera))
		left = vector.Cross3(oc, up)

		invWidth := 1.0 / float64(width)
		invHeight := 1.0 / float64(height)

		for i := 0; i < width; i++ {

			fmt.Println(i)

			iFloat = float64(i)*invWidth - 0.5*(1-invWidth)

			for j := 0; j < height; j++ {

				jFloat = float64(j)*invHeight - 0.5*(1-invHeight)

				dir = vector.Sum3(vector.Sum3(oc, vector.Scale3(up, jFloat)), vector.Scale3(left, iFloat))

				colour = resources.RayTrace(dir, camera, configData.Eta1, [3]float64{configData.Eta2R, configData.Eta2G, configData.Eta2B}, configData.NumberOfBounces, faceData, up, left, invHeight, invWidth, configData.RaysPerPixel)

				r = colour[0]
				g = colour[1]
				b = colour[2]
				a = 255

				r = 255 * r
				g = 255 * g
				b = 255 * b

				col = color.RGBA{
					uint8(r),
					uint8(g),
					uint8(b),
					uint8(a),
				}

				img.Set(i, j, col)

			}

		}

		f, _ := os.Create(fmt.Sprintf("images/png/data%d.png", time))
		png.Encode(f, img)

		fmt.Println(time)

	}

	timeString := time.Now().Format("01-02-2006 15:04:05")

	// Read from the config file to determine properties for the gif
	origJson, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}

	copyJson, err := os.Create("images/data/" + timeString + ".json")
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
