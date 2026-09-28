package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/luanorlando/labs-cep-weather-obsevability.git/internal/entity"
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

func (c WeatherServiceClient) Fetch(cep string) (*WeatherOutputDto, error) {
	if !entity.ValidarCEP(cep) {
		return nil, entity.ErrCEPInvalid
	}

	client := http.Client{}
	urlAPI := fmt.Sprintf("https://viacep.com.br/ws/%s/json", cep)
	resp, err := client.Get(urlAPI)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info entity.Cep
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	if info.Erro == "true" {
		return nil, entity.ErrCEPNotFound
	}

	return c.fetchBy(info.City, &client)
}

func (c WeatherServiceClient) fetchBy(city string, client *http.Client) (*WeatherOutputDto, error) {
	escapedCity := url.QueryEscape(city)
	apiUrl := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?q=%s&lang=pt&key=%s", escapedCity, c.apiKey)

	fmt.Printf("URL API: %s", apiUrl)
	resp, err := client.Get(apiUrl)
	if err != nil {
		fmt.Printf("erro: %s", err.Error())
		return nil, entity.ErrWeatherFound
	}
	defer resp.Body.Close()

	var weather entity.WeatherFromCity

	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		fmt.Printf("Erro decode: %s", err.Error())
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
