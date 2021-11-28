package resources

import (
	"image/color"
	"math"
)

type whitePoint struct {
	xWhite float64
	yWhite float64
}

type colourSystem struct {
	name            string
	xRed            float64
	yRed            float64
	xGreen          float64
	yGreen          float64
	xBlue           float64
	yBlue           float64
	whitePoint      whitePoint
	gammaCorrection float64
}

const GAMMA_REC709 = 0

var (
	IlluminantC   = whitePoint{0.3101, 0.3162}
	IlluminantD65 = whitePoint{0.3127, 0.3291}
	IlluminantE   = whitePoint{0.33333333, 0.33333333}
)

var (
	NTSCsystem   = colourSystem{"NTSC", 0.67, 0.33, 0.21, 0.71, 0.14, 0.08, IlluminantC, GAMMA_REC709}
	EBUsystem    = colourSystem{"EBU (PAL/SECAM)", 0.64, 0.33, 0.29, 0.60, 0.15, 0.06, IlluminantD65, GAMMA_REC709}
	SMPTEsystem  = colourSystem{"SMPTE", 0.630, 0.340, 0.310, 0.595, 0.155, 0.070, IlluminantD65, GAMMA_REC709}
	HDTVsystem   = colourSystem{"HDTV", 0.670, 0.330, 0.210, 0.710, 0.150, 0.060, IlluminantD65, GAMMA_REC709}
	CIEsystem    = colourSystem{"CIE", 0.7355, 0.2645, 0.2658, 0.7243, 0.1669, 0.0085, IlluminantE, GAMMA_REC709}
	Rec709system = colourSystem{"CIE REC 709", 0.64, 0.33, 0.30, 0.60, 0.15, 0.06, IlluminantD65, GAMMA_REC709}
)

var (
	XMatchFunction = []float64{0.0014, 0.0022, 0.0042, 0.0076, 0.0143, 0.0232, 0.0435, 0.0776, 0.1344, 0.2148, 0.2839, 0.3285, 0.3483, 0.3481, 0.3362, 0.3187, 0.2908, 0.2511, 0.1954, 0.1421, 0.0956, 0.058, 0.032, 0.0147, 0.0049, 0.0024, 0.0093, 0.0291, 0.0633, 0.1096, 0.1655, 0.2257, 0.2904, 0.3597, 0.4334, 0.5121, 0.5945, 0.6784, 0.7621, 0.8425, 0.9163, 0.9786, 1.0263, 1.0567, 1.0622, 1.0456, 1.0026, 0.9384, 0.8544, 0.7514, 0.6424, 0.5419, 0.4479, 0.3608, 0.2835, 0.2187, 0.1649, 0.1212, 0.0874, 0.0636, 0.0468, 0.0329, 0.0227, 0.0158, 0.0114, 0.0081, 0.0058, 0.0041, 0.0029, 0.002, 0.0014, 0.001, 0.0007, 0.0005, 0.0003, 0.0002, 0.0002, 0.0001, 0.0001, 0.0001, 0}
	YMatchFunction = []float64{0, 0.0001, 0.0001, 0.0002, 0.0004, 0.0006, 0.0012, 0.0022, 0.004, 0.0073, 0.0116, 0.0168, 0.023, 0.0298, 0.038, 0.048, 0.06, 0.0739, 0.091, 0.1126, 0.139, 0.1693, 0.208, 0.2586, 0.323, 0.4073, 0.503, 0.6082, 0.71, 0.7932, 0.862, 0.9149, 0.954, 0.9803, 0.995, 1, 0.995, 0.9786, 0.952, 0.9154, 0.87, 0.8163, 0.757, 0.6949, 0.631, 0.5668, 0.503, 0.4412, 0.381, 0.321, 0.265, 0.217, 0.175, 0.1382, 0.107, 0.0816, 0.061, 0.0446, 0.032, 0.0232, 0.017, 0.0119, 0.0082, 0.0057, 0.0041, 0.0029, 0.0021, 0.0015, 0.001, 0.0007, 0.0005, 0.0004, 0.0002, 0.0002, 0.0001, 0.0001, 0.0001, 0, 0, 0, 0}
	ZMatchFunction = []float64{0.0065, 0.0105, 0.0201, 0.0362, 0.0679, 0.1102, 0.2074, 0.3713, 0.6456, 1.0391, 1.3856, 1.623, 1.7471, 1.7826, 1.7721, 1.7441, 1.6692, 1.5281, 1.2876, 1.0419, 0.813, 0.6162, 0.4652, 0.3533, 0.272, 0.2123, 0.1582, 0.1117, 0.0782, 0.0573, 0.0422, 0.0298, 0.0203, 0.0134, 0.0087, 0.0057, 0.0039, 0.0027, 0.0021, 0.0018, 0.0017, 0.0014, 0.0011, 0.001, 0.0008, 0.0006, 0.0003, 0.0002, 0.0002, 0.0001, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

var Y_white = 0.0

func Xyztorgb(cs colourSystem, xc, yc, zc float64, scale bool) (r, g, b float64) {

	if scale {
		xc /= Y_white
		yc /= Y_white
		zc /= Y_white
	}

	var xr, yr, zr, xg, yg, zg, xb, yb, zb float64
	var xw, yw, zw float64
	var rx, ry, rz, gx, gy, gz, bx, by, bz float64
	var rw, gw, bw float64

	xr = cs.xRed
	yr = cs.yRed
	zr = 1 - (xr + yr)
	xg = cs.xGreen
	yg = cs.yGreen
	zg = 1 - (xg + yg)
	xb = cs.xBlue
	yb = cs.yBlue
	zb = 1 - (xb + yb)

	xw = cs.whitePoint.xWhite
	yw = cs.whitePoint.yWhite
	zw = 1 - (xw + yw)

	rx = (yg * zb) - (yb * zg)
	ry = (xb * zg) - (xg * zb)
	rz = (xg * yb) - (xb * yg)
	gx = (yb * zr) - (yr * zb)
	gy = (xr * zb) - (xb * zr)
	gz = (xb * yr) - (xr * yb)
	bx = (yr * zg) - (yg * zr)
	by = (xg * zr) - (xr * zg)
	bz = (xr * yg) - (xg * yr)

	rw = ((rx * xw) + (ry * yw) + (rz * zw))
	gw = ((gx * xw) + (gy * yw) + (gz * zw))
	bw = ((bx * xw) + (by * yw) + (bz * zw))

	rx = rx / rw
	ry = ry / rw
	rz = rz / rw
	gx = gx / gw
	gy = gy / gw
	gz = gz / gw
	bx = bx / bw
	by = by / bw
	bz = bz / bw

	r = (rx * xc) + (ry * yc) + (rz * zc)
	g = (gx * xc) + (gy * yc) + (gz * zc)
	b = (bx * xc) + (by * yc) + (bz * zc)

	return r, g, b
}

func Constrainrgb(r, g, b float64) (float64, float64, float64) {

	w := math.Min(r, math.Min(g, b))

	if w > 1.0 {

		ww := math.Max(r, math.Max(g, b)) - 1

		r -= ww
		g -= ww
		b -= ww

	}

	w = math.Max(r, math.Max(g, b))

	for w < 0.0 {

		r += 1.0
		g += 1.0
		b += 1.0

		w = math.Max(r, math.Max(g, b))

	}

	return r, g, b

}

func GammaCorrect(c float64) float64 {

	if c > 0.0031308 {
		return 1.055*math.Pow(c, 0.45) - 0.055
	} else {
		return 12.92 * c
	}

}

func Gammacorrect(rgb [3]float64) [3]float64 {

	return [3]float64{GammaCorrect(rgb[0]), GammaCorrect(rgb[1]), GammaCorrect(rgb[2])}

}

func Normrgb(x [3]float64) [3]float64 {

	rgb := [3]float64{0, 0, 0}
	greatest := math.Max(x[0], math.Max(x[1], x[2]))

	if greatest > 0 {
		rgb[0] = x[0] / greatest
		rgb[1] = x[1] / greatest
		rgb[2] = x[2] / greatest
	}

	return rgb

}

func IntegrateSpectrum(spectrum []float64, scale float64) [3]float64 {

	x := 0.0
	y := 0.0
	z := 0.0

	for i := 0; i < len(spectrum); i++ {

		x += XMatchFunction[i] * spectrum[i]
		y += YMatchFunction[i] * spectrum[i]
		z += ZMatchFunction[i] * spectrum[i]

	}

	return [3]float64{x / scale, y / scale, z / scale}

}

func XYZToRGB(x, y, z float64) (float64, float64, float64) {

	r := math.Min(math.Max(3.2406*x-1.5372*y-0.4986*z, 0), 1)
	g := math.Min(math.Max(-0.9689*x+1.8758*y+0.0415*z, 0), 1)
	b := math.Min(math.Max(0.0557*x-0.2040*y+1.0570*z, 0), 1)

	return r, g, b

}

func SpectrumToRGBA(spectrum []float64, scale float64, sigma float64) color.RGBA {

	spectrum = blur(spectrum, sigma)
	xyz := IntegrateSpectrum(spectrum, scale)

	r, g, b := Constrainrgb(xyz[0], xyz[1], xyz[2])

	r, g, b = XYZToRGB(r, g, b)

	r = GammaCorrect(r)
	g = GammaCorrect(g)
	b = GammaCorrect(b)

	r *= 255.0
	g *= 255.0
	b *= 255.0

	return color.RGBA{
		uint8(r),
		uint8(g),
		uint8(b),
		255,
	}

}

func blur(spectrum []float64, sigma float64) []float64 {

	blurredSpectrum := make([]float64, len(spectrum))

	for i := 0; i < len(spectrum); i++ {

		if spectrum[i] > 0.0 {

			for j := 0; j < len(spectrum); j++ {

				if float64(i-j) < sigma && float64(i-j) > -sigma {

					blurredSpectrum[j] += spectrum[i] * math.Exp(-float64(i-j)*float64(i-j)/(2*sigma*sigma)) / (sigma * math.Sqrt(2*math.Pi))

				}

			}

		}

	}

	return blurredSpectrum

}

func MergeColourStimulus(spectralRaysNumber int) {

	newX := make([]float64, spectralRaysNumber)
	newY := make([]float64, spectralRaysNumber)
	newZ := make([]float64, spectralRaysNumber)

	i := 0
	for i < 80 {

		newX[i*spectralRaysNumber/80] += XMatchFunction[i]
		newY[i*spectralRaysNumber/80] += YMatchFunction[i]
		newZ[i*spectralRaysNumber/80] += ZMatchFunction[i]

		i++

	}

	for i := 0; i < spectralRaysNumber; i++ {

		newX[i] *= float64(spectralRaysNumber) / 80.0
		newY[i] *= float64(spectralRaysNumber) / 80.0
		newZ[i] *= float64(spectralRaysNumber) / 80.0

	}

	XMatchFunction = newX
	YMatchFunction = newY
	ZMatchFunction = newZ
}
