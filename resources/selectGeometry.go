package resources

import "github.com/calummccain/coxeter/data"

func SelectGeometry(p, q int, r float64, truncRect string, numberOfFaces int) data.CellData {

	var cellData data.CellData
	cell := [2]int{p, q}

	if truncRect == "r" {

		switch cell {
		case [2]int{3, 3}:
			cellData = data.RectifiedTetrahedronData(r)
		case [2]int{3, 4}:
			cellData = data.RectifiedOctahedronData(r)
		case [2]int{3, 5}:
			cellData = data.RectifiedIcosahedronData(r)
		case [2]int{4, 3}:
			cellData = data.RectifiedHexahedronData(r)
		case [2]int{5, 3}:
			cellData = data.RectifiedDodecahedronData(r)
		default:
			cellData = data.TetrahedronData(r)
		}

	} else if truncRect == "t" {

		switch cell {
		case [2]int{3, 3}:
			cellData = data.TruncatedTetrahedronData(r)
		case [2]int{3, 4}:
			cellData = data.TruncatedOctahedronData(r)
		case [2]int{3, 5}:
			cellData = data.TruncatedIcosahedronData(r)
		case [2]int{4, 3}:
			cellData = data.TruncatedHexahedronData(r)
		case [2]int{5, 3}:
			cellData = data.TruncatedDodecahedronData(r)
		default:
			cellData = data.TetrahedronData(r)
		}

	} else {

		switch cell {
		case [2]int{3, 3}:
			cellData = data.TetrahedronData(r)
		case [2]int{3, 4}:
			cellData = data.OctahedronData(r)
		case [2]int{3, 5}:
			cellData = data.IcosahedronData(r)
		case [2]int{3, 6}:
			cellData = data.TriangularData(r, numberOfFaces)
			cellData.C = [4]float64{1, -4, 0, 0}
		case [2]int{4, 3}:
			cellData = data.HexahedronData(r)
		case [2]int{4, 4}:
			cellData = data.SquareData(r, numberOfFaces)
		case [2]int{5, 3}:
			cellData = data.DodecahedronData(r)
		case [2]int{6, 3}:
			cellData = data.HexagonalData(r, numberOfFaces)
		default:
			cellData = data.HyperbolicData(p, q, r, numberOfFaces)
		}

	}

	return cellData

}
