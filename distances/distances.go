package distances

import (
	"math"
)

const (
	// The Earth's mean radius in meters (according to NASA).
	earthRadiusm     = 6371
	maxGraphDistance = math.MaxFloat32
)

// CalculateDistanceMeters calculate lineal distance in meters between two points using Manhattan formula.
func CalculateDistanceMeters(ori, dest [2]float64) float64 {
	x1 := ori[0] - dest[0]
	dLat := x1 * math.Pi / 180
	x2 := ori[1] - dest[1]
	dLon := x2 * math.Pi / 180
	latDistanceA := math.Sin(dLat/2) * math.Sin(dLat/2)
	latDistanceB := 2 * math.Atan2(math.Sqrt(latDistanceA), math.Sqrt(1-latDistanceA))
	latDistance := earthRadiusm * latDistanceB
	longDistanceA := math.Sin(dLon/2) * math.Sin(dLon/2)
	longDistanceB := 2 * math.Atan2(math.Sqrt(longDistanceA), math.Sqrt(1-longDistanceA))
	longDistance := earthRadiusm * longDistanceB
	distOutput := math.Abs(latDistance) + math.Abs(longDistance)
	return distOutput * 1000
}

func IsInvalidDistance(d float32) bool {
	return d >= maxGraphDistance
}

func MetersToKm(meters float64) float64 {
	return meters / 1000
}
