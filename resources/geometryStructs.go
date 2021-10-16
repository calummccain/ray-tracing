package resources

type VertexEuclidean struct {
	V4 [4]float64
	V3 [3]float64
}

type VertexSpherical struct {
	V4 [4]float64
	V3 [3]float64
}

type VertexHyperbolic struct {
	V4 [4]float64
	V3 [3]float64
}

type FaceEuclidean struct {
	D      float64
	Normal [3]float64
}

type FaceSpherical struct {
	Type         string
	D            float64
	Normal       [3]float64
	Radius       float64
	SphereCenter [3]float64
	Center3      [3]float64
	Center4      [4]float64
	Polygon3     [][3]float64
	Polygon4     [][4]float64
	Plane4       [][4]float64
	Dot4         []float64
}

type FaceHyperbolic struct {
	Type         string
	Radius       float64
	SphereCenter [3]float64
	D            float64
	Normal       [3]float64
	Plane4       [][4]float64
	Dot4         []float64
	InOut        bool
}

type PointVisibility struct {
	Point   [3]float64
	Visible bool
}
