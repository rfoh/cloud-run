package main

import (
	"log"
	"net/http"
	"os"

	"rfoh/cloud-run/internal/adapter/external/viacep"
	"rfoh/cloud-run/internal/adapter/external/weatherapi"
	httpAdapter "rfoh/cloud-run/internal/adapter/http"
	"rfoh/cloud-run/internal/application/usecase"
)

func main() {
	httpClient := &http.Client{}

	cepProvider := viacep.NewViaCEPAdapter(httpClient)
	weatherProvider := weatherapi.NewWeatherAPIAdapter(httpClient)

	getTemperatureUseCase := usecase.NewGetTemperatureByCEPUseCase(cepProvider, weatherProvider)

	weatherHandler := httpAdapter.NewWeatherHandler(getTemperatureUseCase)

	http.HandleFunc("/", weatherHandler.Handle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
