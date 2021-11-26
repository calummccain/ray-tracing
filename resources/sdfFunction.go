package resources

func SdfFunction(config Config, faceData [][]Face) (sdf func([3]float64) float64) {

	if config.ObjectConfig.Sdf == "seh" {

		sdf = func(p [3]float64) float64 {
			return Sdf(p, faceData, config.CellGeometryConfig.Model)
		}

	} else {

		if config.ObjectConfig.Sdf == "sphere" {

			sdf = func(p [3]float64) float64 {
				return SdfSphere(RotateXYZ(p, config.ObjectConfig.ObjectRotateX, config.ObjectConfig.ObjectRotateY, config.ObjectConfig.ObjectRotateZ), config.SphereConfig)
			}

		} else if config.ObjectConfig.Sdf == "cube" {

			sdf = func(p [3]float64) float64 {
				return SdfCube(RotateXYZ(p, config.ObjectConfig.ObjectRotateX, config.ObjectConfig.ObjectRotateY, config.ObjectConfig.ObjectRotateZ), config.CubeConfig)
			}

		} else if config.ObjectConfig.Sdf == "torus" {

			sdf = func(p [3]float64) float64 {
				return SdfTorus(RotateXYZ(p, config.ObjectConfig.ObjectRotateX, config.ObjectConfig.ObjectRotateY, config.ObjectConfig.ObjectRotateZ), config.TorusConfig)
			}

		} else if config.ObjectConfig.Sdf == "spheres" {

			sdf = func(p [3]float64) float64 {
				return SdfSpheres(RotateXYZ(p, config.ObjectConfig.ObjectRotateX, config.ObjectConfig.ObjectRotateY, config.ObjectConfig.ObjectRotateZ))
			}

		}

	}

	return sdf

}
