package resources

import (
	"github.com/calummccain/coxeter/data"
)

func SelectGeometry(CGC CellGeometryConfig) data.CellData {

	var cellData data.CellData
	//cell := [2]int{CGC.P, CGC.Q}

	if CGC.TruncRect == "r" {

		switch [2]int{CGC.P, CGC.Q} {
		case [2]int{3, 3}:
			cellData = data.RectifiedTetrahedronData(CGC.R)
		case [2]int{3, 4}:
			cellData = data.RectifiedOctahedronData(CGC.R)
		case [2]int{3, 5}:
			cellData = data.RectifiedIcosahedronData(CGC.R)
		case [2]int{4, 3}:
			cellData = data.RectifiedHexahedronData(CGC.R)
		case [2]int{5, 3}:
			cellData = data.RectifiedDodecahedronData(CGC.R)
		default:
			cellData = data.TetrahedronData(CGC.R)
		}

	} else if CGC.TruncRect == "t" {

		switch [2]int{CGC.P, CGC.Q} {
		case [2]int{3, 3}:
			cellData = data.TruncatedTetrahedronData(CGC.R)
		case [2]int{3, 4}:
			cellData = data.TruncatedOctahedronData(CGC.R)
		case [2]int{3, 5}:
			cellData = data.TruncatedIcosahedronData(CGC.R)
		case [2]int{4, 3}:
			cellData = data.TruncatedHexahedronData(CGC.R)
		case [2]int{5, 3}:
			cellData = data.TruncatedDodecahedronData(CGC.R)
		default:
			cellData = data.TetrahedronData(CGC.R)
		}

	} else {

		switch [2]int{CGC.P, CGC.Q} {
		case [2]int{3, 3}:
			cellData = data.TetrahedronData(CGC.R)
		case [2]int{3, 4}:
			cellData = data.OctahedronData(CGC.R)
		case [2]int{3, 5}:
			cellData = data.IcosahedronData(CGC.R)
		case [2]int{3, 6}:
			cellData = data.TriangularData(CGC.R, CGC.NumberOfFaces)
			cellData.C = [4]float64{1, -4, 0, 0}
		case [2]int{4, 3}:
			cellData = data.HexahedronData(CGC.R)
		case [2]int{4, 4}:
			cellData = data.SquareData(CGC.R, CGC.NumberOfFaces)
		case [2]int{5, 3}:
			cellData = data.DodecahedronData(CGC.R)
		case [2]int{6, 3}:
			cellData = data.HexagonalData(CGC.R, CGC.NumberOfFaces)
		default:
			cellData = data.HyperbolicData(CGC.P, CGC.Q, CGC.R, CGC.NumberOfFaces)
		}

	}

	return cellData

}
