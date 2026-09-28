package usecase

import "github.com/luanorlando/labs-cep-weather-obsevability.git/service-b/internal/repository"

type WeatherServiceClient interface {
	Fetch(cep string) (*repository.WeatherOutputDto, error)
}

type FetchWeather struct {
	client WeatherServiceClient
}

func NewFetchWeatherUsecase(c WeatherServiceClient) *FetchWeather {
	return &FetchWeather{
		client: c,
	}
}

func (u FetchWeather) Execute(cep string) (*repository.WeatherOutputDto, error) {
	result, err := u.client.Fetch(cep)
	if err != nil {
		return nil, err
	}

	return result, nil
}
