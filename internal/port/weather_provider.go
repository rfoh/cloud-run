package port

type WeatherProvider interface {
	GetTemperature(city string) (float64, error)
}
