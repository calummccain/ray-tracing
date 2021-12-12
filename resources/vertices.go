package resources

import (
	"github.com/calummccain/coxeter/hyperbolic"
	"github.com/calummccain/coxeter/shared"
	"github.com/calummccain/coxeter/spherical"
	"github.com/calummccain/coxeter/vector"
)

func GenerateVerticesSpherical(vertices [][4]float64, cell string, matrices shared.Matrices, numVertices int) []VertexSpherical {

	newVertices := vector.TransformVertices(vertices, cell, matrices)
	var verts []VertexSpherical
	var p [4]float64

	for i := 0; i < numVertices; i++ {

		p = matrices.F(newVertices[i])

		verts = append(verts, VertexSpherical{
			V4: p,
			V3: spherical.HyperToStereo(p),
		})

	}

	return verts

}

func GenerateVerticesEuclidean(vertices [][4]float64, cell string, matrices shared.Matrices, numVertices int) []VertexEuclidean {

	newVertices := vector.TransformVertices(vertices, cell, matrices)
	var verts []VertexEuclidean
	var p [4]float64

	for i := 0; i < numVertices; i++ {

		p = matrices.F(newVertices[i])

		verts = append(verts, VertexEuclidean{
			V4: p,
			V3: [3]float64{p[1], p[2], p[3]},
		})

	}

	return verts

}

func GenerateVerticesHyperbolic(vertices [][4]float64, cell string, matrices shared.Matrices, numVertices int) []VertexHyperbolic {

	newVertices := vector.TransformVertices(vertices, cell, matrices)
	var verts []VertexHyperbolic
	var p [4]float64

	for i := 0; i < numVertices; i++ {

		p = matrices.F(newVertices[i])

		verts = append(verts, VertexHyperbolic{
			V4: p,
			V3: hyperbolic.HyperboloidToKlein(p),
		})

	}

	return verts

}
