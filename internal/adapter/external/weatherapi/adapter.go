package weatherapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"rfoh/cloud-run/internal/domain/entity"
)

const weatherAPIURL = "https://api.weatherapi.com/v1/current.json"

type weatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type WeatherAPIAdapter struct {
	client *http.Client
	apiKey string
}

func NewWeatherAPIAdapter(client *http.Client) *WeatherAPIAdapter {
	return &WeatherAPIAdapter{
		client: client,
		apiKey: os.Getenv("WEATHER_API_KEY"),
	}
}

func (a *WeatherAPIAdapter) GetTemperature(city string) (float64, error) {
	if a.apiKey == "" {
		return 0, &entity.TemperatureError{City: city, Message: "WEATHER_API_KEY não configurada"}
	}

	params := url.Values{}
	params.Add("key", a.apiKey)
	params.Add("q", city)
	params.Add("aqi", "no")

	resp, err := a.client.Get(weatherAPIURL + "?" + params.Encode())
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, &entity.TemperatureError{City: city, Message: "erro ao buscar temperatura"}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var data weatherAPIResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}

	return data.Current.TempC, nil
}
