package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"rfoh/cloud-run/internal/application/dto"
	"rfoh/cloud-run/internal/domain/entity"
)

type errorResponse struct {
	Message string `json:"message"`
}

type GetTemperatureByCEPUseCaseInterface interface {
	Execute(input *dto.GetTemperatureByCEPInput) (*dto.GetTemperatureByCEPOutput, error)
}

type WeatherHandler struct {
	usecase GetTemperatureByCEPUseCaseInterface
}

func NewWeatherHandler(usecase GetTemperatureByCEPUseCaseInterface) *WeatherHandler {
	return &WeatherHandler{usecase: usecase}
}

func (h *WeatherHandler) Handle(w http.ResponseWriter, r *http.Request) {
	cep := strings.TrimSpace(r.URL.Query().Get("cep"))

	if cep == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Message: "missing cep parameter"})
		return
	}

	input := &dto.GetTemperatureByCEPInput{CEP: cep}
	output, err := h.usecase.Execute(input)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		if errors.Is(err, &entity.CEPNotFoundError{}) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errorResponse{Message: "can not find zipcode"})
		} else if errors.Is(err, &entity.InvalidCEPError{}) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(errorResponse{Message: "invalid zipcode"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{Message: "unknown error"})
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
