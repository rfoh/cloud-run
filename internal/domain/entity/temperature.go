package entity

import "math"

type Temperature struct {
	Celsius    float64
	Fahrenheit float64
	Kelvin     float64
}

func NewTemperature(tempC float64) *Temperature {
	tempF := tempC*1.8 + 32
	tempK := tempC + 273

	return &Temperature{
		Celsius:    round(tempC, 1),
		Fahrenheit: round(tempF, 1),
		Kelvin:     round(tempK, 1),
	}
}

func round(value float64, decimals int) float64 {
	multiplier := math.Pow(10, float64(decimals))
	return math.Round(value*multiplier) / multiplier
}
