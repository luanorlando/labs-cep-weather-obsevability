package usecase

import (
	"context"

	"github.com/luanorlando/labs-cep-weather-obsevability.git/service-b/internal/repository"
)

type WeatherServiceClient interface {
	Fetch(ctx context.Context, cep string) (*repository.WeatherOutputDto, error)
}

type FetchWeather struct {
	client WeatherServiceClient
}

func NewFetchWeatherUsecase(c WeatherServiceClient) *FetchWeather {
	return &FetchWeather{
		client: c,
	}
}

func (u FetchWeather) Execute(ctx context.Context, cep string) (*repository.WeatherOutputDto, error) {
	result, err := u.client.Fetch(ctx, cep)
	if err != nil {
		return nil, err
	}

	return result, nil
}
