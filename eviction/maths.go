package eviction

import "math"

func easeScaler(x float64) float64 {
	proportionalScale := math.Pow(2, 29)
	proportionalScaler := x / proportionalScale

	centeredEase := math.Atan(proportionalScaler - 11)
	scaledEase := centeredEase / math.Pi

	scaleBias := 0.53523

	return scaleBias - scaledEase
}

func scoreData(popularity, bytes int, ageInSeconds float64) float64 {
	if popularity == 0 || bytes == 0 || ageInSeconds == 0 {
		return 0
	}

	popularityVector := math.Pow(2.4, float64(popularity-1))
	popVecLimit := math.Pow(2, 24)
	boundedPopVec := math.Min(popularityVector, popVecLimit)

	sizeScaler := math.Log(float64(bytes)) / math.Log(256)
	sizeScalerDelta := math.Log(math.Pow(2, 22)) / math.Log(256)
	sizeScaledVector := sizeScaler - sizeScalerDelta

	sizeDrvOrientation := 0.0
	if sizeScaledVector > 0 {
		sizeDrvOrientation = 1.0
	}
	if sizeScaledVector < 0 {
		sizeDrvOrientation = -1.0
	}

	sizeAsymptote := sizeScaledVector * sizeDrvOrientation
	sizeVector := 1 / (sizeAsymptote + 1)

	ageScaler := ageInSeconds / (60 * 60 * 24)
	ageVector := math.Pow(6, ageScaler)

	popVecScaler := sizeVector * ageVector
	omniScaler := easeScaler(float64(bytes))
	omniScalerBounds := easeScaler(math.Pow(2, 30))
	boundedOmniScaler := math.Min(omniScaler, omniScalerBounds)

	return (boundedPopVec / popVecScaler) * boundedOmniScaler
}
