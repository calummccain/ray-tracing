package resources

import "image/color"

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

type Face struct {
	Type         string
	Radius       float64
	SphereCenter [3]float64
	D            float64
	Normal       [3]float64
	Plane4       [][4]float64
	Dot4         []float64
	Outside      bool
}

type PointVisibility struct {
	Point   [3]float64
	Visible bool
}

type ray struct {
	pos     [3]float64
	dir     [3]float64
	weight  float64
	rayType string
	parent  int
	inside  bool
	layer   int
	length  float64
}

type Light struct {
	Temp      float64
	Intensity float64
	Spectrum  []float64
	Pos       [3]float64
	Up        [3]float64
	Left      [3]float64
	Normal    [3]float64
	Height    float64
	Width     float64
	Inside    bool
}

type ColumnColours struct {
	ColumnNumber int
	Colours      []color.RGBA
}

func (l *Light) SpectrumFromWavelength(wavelengths []float64, SpectralRaysNumber int) {

	for _, lambda := range wavelengths {

		l.Spectrum = append(l.Spectrum, l.Intensity*(5/float64(SpectralRaysNumber))*BlackBodySpectrum(lambda, l.Temp))

	}

}

func (l *Light) LightInside(sdf func([3]float64) float64) {

	if sdf(l.Pos) < 0.0 {
		l.Inside = true
	} else {
		l.Inside = false
	}

}
