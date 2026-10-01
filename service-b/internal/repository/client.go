package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/luanorlando/labs-cep-weather-obsevability.git/internal/entity"
	"go.opentelemetry.io/otel"
)

type WeatherOutputDto struct {
	City       string  `json:"city"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_k"`
}

type WeatherServiceClient struct {
	apiKey string
}

func NewWeatherServiceClient(apiKey string) *WeatherServiceClient {
	return &WeatherServiceClient{
		apiKey: apiKey,
	}
}

func (c WeatherServiceClient) Fetch(ctx context.Context, cep string) (*WeatherOutputDto, error) {
	if !entity.ValidarCEP(cep) {
		return nil, entity.ErrCEPInvalid
	}

	tracer := otel.Tracer("service-b-repository")
	cepCtx, spam := tracer.Start(ctx, "FetchExternalCEPAPI")
	defer spam.End()

	client := http.Client{}
	urlAPI := fmt.Sprintf("https://viacep.com.br/ws/%s/json", cep)
	weatherReq, err := http.NewRequestWithContext(cepCtx, "GET", urlAPI, nil)
	if err != nil {
		spam.RecordError(err)
		return nil, entity.ErrCEPNotFound
	}

	resp, err := client.Do(weatherReq)

	if err != nil {
		spam.RecordError(err)
		return nil, entity.ErrCEPNotFound
	}
	defer resp.Body.Close()

	var info entity.Cep
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		spam.RecordError(err)
		return nil, err
	}

	if info.Erro == "true" {
		spam.RecordError(entity.ErrCEPNotFound)
		return nil, entity.ErrCEPNotFound
	}

	return c.fetchBy(cepCtx, info.City, &client)
}

func (c WeatherServiceClient) fetchBy(ctx context.Context, city string, client *http.Client) (*WeatherOutputDto, error) {
	escapedCity := url.QueryEscape(city)
	apiUrl := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?q=%s&lang=pt&key=%s", escapedCity, c.apiKey)

	tracer := otel.Tracer("service-b-repository")
	weatherCtx, spam := tracer.Start(ctx, "FetchExternalWeatherAPI")
	weatherReq, err := http.NewRequestWithContext(weatherCtx, "GET", apiUrl, nil)

	resp, err := client.Do(weatherReq)

	if err != nil {
		spam.RecordError(err)
		return nil, entity.ErrWeatherFound
	}
	defer resp.Body.Close()

	var weather entity.WeatherFromCity

	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		decodeErr := fmt.Errorf("Erro decode: %s", err.Error())
		spam.RecordError(decodeErr)
		return nil, err
	}

	fmt.Printf("Response:")
	fmt.Printf("A temperatura atual é: %.1f°C para a cidade de %s", weather.Current.Celsius, city)

	kelvin := entity.KelvinBy(weather.Current.Celsius)

	return &WeatherOutputDto{
		City:       city,
		Celsius:    weather.Current.Celsius,
		Fahrenheit: weather.Current.Fahrenheit,
		Kelvin:     kelvin,
	}, nil
}
