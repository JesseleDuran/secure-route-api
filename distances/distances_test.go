package distances

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CalculateDistance(t *testing.T) {
	var tt = []struct {
		origin       [2]float64
		destination  [2]float64
		outputMeters float64
	}{
		{
			[2]float64{-32.93196245876229, -71.54468289258033},  // Concon, Chile
			[2]float64{-32.971730342514114, -71.54147365771965}, // Viña del Mar, Chile
			4778.837551511225,
		},
		{
			[2]float64{-32.971730342514114, -71.54147365771965}, // Viña del Mar, Chile
			[2]float64{-33.42625728925902, -70.60862030177017},  // Santiago, Chile
			154269.6509862075,
		},
		{
			[2]float64{10.475386530483375, -66.79629254579477}, // Caracas, Venezuela
			[2]float64{10.33077132679323, -67.08100156128076},  // Los Teques, Venezuela
			47738.67505802252,
		},
	}
	for _, input := range tt {
		got := CalculateDistanceMeters(input.origin, input.destination)

		assert.NotNil(t, got)
		assert.Equal(t, input.outputMeters, got)
	}
}

func TestMetersToKm(t *testing.T) {
	var tt = []struct {
		meters           float64
		resultKilometers float64
	}{
		{
			12345,
			12.345,
		},
		{
			-50,
			-0.05,
		},
		{
			0,
			0,
		},
	}
	for _, input := range tt {
		got := MetersToKm(input.meters)

		assert.NotNil(t, got)
		assert.Equal(t, input.resultKilometers, got)
	}
}
