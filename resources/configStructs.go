package resources

type CellGeometryData struct {
	P             int
	Q             int
	R             float64
	Cells         []string
	TruncRect     string
	Model         string
	NumberOfFaces int
}

type Config struct {
	CellGeometryData CellGeometryData
	Width            int
	Height           int
	Start            int
	End              int
	Sdf              string
	Distance         float64
	Eta1             float64
	Eta2R            float64
	Eta2G            float64
	Eta2B            float64
	NumberOfBounces  int
	RaysPerPixel     int
	SphereRadius     float64
	CubeA            float64
	CubeB            float64
	CubeC            float64
	TorusA           float64
	TorusB           float64
	ObjectRotateX    float64
	ObjectRotateY    float64
	ObjectRotateZ    float64
	CameraRotateX    float64
	CameraRotateY    float64
	CameraRotateZ    float64
	Save             bool
	SpectralRays     int
	Temp             int
	Spectral         bool
}
