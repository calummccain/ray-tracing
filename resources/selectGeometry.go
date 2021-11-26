package resources

import "github.com/calummccain/coxeter/data"

func SelectGeometry(CGD CellGeometryData) data.CellData {

	var cellData data.CellData
	cell := [2]int{CGD.P, CGD.Q}

	if CGD.TruncRect == "r" {

		switch cell {
		case [2]int{3, 3}:
			cellData = data.RectifiedTetrahedronData(CGD.R)
		case [2]int{3, 4}:
			cellData = data.RectifiedOctahedronData(CGD.R)
		case [2]int{3, 5}:
			cellData = data.RectifiedIcosahedronData(CGD.R)
		case [2]int{4, 3}:
			cellData = data.RectifiedHexahedronData(CGD.R)
		case [2]int{5, 3}:
			cellData = data.RectifiedDodecahedronData(CGD.R)
		default:
			cellData = data.TetrahedronData(CGD.R)
		}

	} else if CGD.TruncRect == "t" {

		switch cell {
		case [2]int{3, 3}:
			cellData = data.TruncatedTetrahedronData(CGD.R)
		case [2]int{3, 4}:
			cellData = data.TruncatedOctahedronData(CGD.R)
		case [2]int{3, 5}:
			cellData = data.TruncatedIcosahedronData(CGD.R)
		case [2]int{4, 3}:
			cellData = data.TruncatedHexahedronData(CGD.R)
		case [2]int{5, 3}:
			cellData = data.TruncatedDodecahedronData(CGD.R)
		default:
			cellData = data.TetrahedronData(CGD.R)
		}

	} else {

		switch cell {
		case [2]int{3, 3}:
			cellData = data.TetrahedronData(CGD.R)
		case [2]int{3, 4}:
			cellData = data.OctahedronData(CGD.R)
		case [2]int{3, 5}:
			cellData = data.IcosahedronData(CGD.R)
		case [2]int{3, 6}:
			cellData = data.TriangularData(CGD.R, CGD.NumberOfFaces)
			cellData.C = [4]float64{1, -4, 0, 0}
		case [2]int{4, 3}:
			cellData = data.HexahedronData(CGD.R)
		case [2]int{4, 4}:
			cellData = data.SquareData(CGD.R, CGD.NumberOfFaces)
		case [2]int{5, 3}:
			cellData = data.DodecahedronData(CGD.R)
		case [2]int{6, 3}:
			cellData = data.HexagonalData(CGD.R, CGD.NumberOfFaces)
		default:
			cellData = data.HyperbolicData(CGD.P, CGD.Q, CGD.R, CGD.NumberOfFaces)
		}

	}

	return cellData

}
