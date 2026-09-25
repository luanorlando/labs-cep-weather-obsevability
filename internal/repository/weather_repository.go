package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/luanorlando/labs-cep-weather-obsevability.git/internal/entity"
)

type ExtenalWeatherRepository struct {
	apiKey     string
	httpClient *http.Client
}

func NewExternalWeatherRepository(apiKey string, client *http.Client) *ExtenalWeatherRepository {
	return &ExtenalWeatherRepository{
		apiKey:     apiKey,
		httpClient: client,
	}
}

func (r ExtenalWeatherRepository) FetchBy(city string) (*entity.WeatherFromCity, error) {
	escapedCity := url.QueryEscape(city)
	apiUrl := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?q=%s&lang=pt&key=%s", escapedCity, r.apiKey)

	fmt.Printf("URL API: %s", apiUrl)
	resp, err := r.httpClient.Get(apiUrl)
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
	fmt.Printf("A temperatura atual é: %.1f°C\n", weather.Current.Celsius)
	fmt.Printf("A temperatura atual é: %.1f°C\n", weather.Current.Fahrenheit)

	return &weather, nil
}
