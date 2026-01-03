package usecase

import (
	"rfoh/cloud-run/internal/application/dto"
	"rfoh/cloud-run/internal/domain/entity"
	"rfoh/cloud-run/internal/port"
)

type GetTemperatureByCEPUseCase struct {
	cepProvider     port.CEPProvider
	weatherProvider port.WeatherProvider
}

func NewGetTemperatureByCEPUseCase(cepProvider port.CEPProvider, weatherProvider port.WeatherProvider) *GetTemperatureByCEPUseCase {
	return &GetTemperatureByCEPUseCase{
		cepProvider:     cepProvider,
		weatherProvider: weatherProvider,
	}
}

func (u *GetTemperatureByCEPUseCase) Execute(input *dto.GetTemperatureByCEPInput) (*dto.GetTemperatureByCEPOutput, error) {
	cep, err := entity.NewCEP(input.CEP)
	if err != nil {
		return nil, err
	}

	location, err := u.cepProvider.FindLocationByCEP(cep)
	if err != nil {
		return nil, err
	}

	tempC, err := u.weatherProvider.GetTemperature(location.City)
	if err != nil {
		return nil, err
	}

	temperature := entity.NewTemperature(tempC)

	return &dto.GetTemperatureByCEPOutput{
		TempC: temperature.Celsius,
		TempF: temperature.Fahrenheit,
		TempK: temperature.Kelvin,
	}, nil
}
