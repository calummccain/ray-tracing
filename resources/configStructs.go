package resources

type CellGeometryConfig struct {
	P             int
	Q             int
	R             float64
	Cells         []string
	TruncRect     string
	Model         string
	NumberOfFaces int
}

type CubeConfig struct {
	CubeA float64
	CubeB float64
	CubeC float64
}

type TorusConfig struct {
	TorusA float64
	TorusB float64
}

type SphereConfig struct {
	SphereRadius float64
}

type ObjectConfig struct {
	Sdf           string
	ObjectRotateX float64
	ObjectRotateY float64
	ObjectRotateZ float64
}

type CameraConfig struct {
	Distance      float64
	CameraRotateX float64
	CameraRotateY float64
	CameraRotateZ float64
}

type ImageConfig struct {
	Width  int
	Height int
	Start  int
	End    int
	Save   bool
}

type RaytracingConfig struct {
	NumberOfBounces    int
	RaysPerPixel       int
	SpectralRaysNumber int
	Spectral           bool
}

type MaterialConfig struct {
	Eta1  float64
	Eta2R float64
	Eta2G float64
	Eta2B float64
}

type LightConfig struct {
	Temp   float64
	Pos    [3]float64
	Up     [3]float64
	Left   [3]float64
	Normal [3]float64
	Height float64
	Width  float64
}

type Config struct {
	CellGeometryConfig CellGeometryConfig
	ObjectConfig       ObjectConfig
	CameraConfig       CameraConfig
	ImageConfig        ImageConfig
	RaytracingConfig   RaytracingConfig
	MaterialConfig     MaterialConfig
	SphereConfig       SphereConfig
	CubeConfig         CubeConfig
	TorusConfig        TorusConfig
	LightsConfig       []LightConfig
}
