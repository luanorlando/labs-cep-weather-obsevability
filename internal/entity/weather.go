package entity

import (
	"errors"
	"regexp"
)

type Weather struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_k"`
}

type Cep struct {
	City string `json:"localidade"`
	Erro string `json:"erro,omitempty"`
}

type WeatherFromCity struct {
	Current CurrentWeather `json:"current"`
}

type CurrentWeather struct {
	Celsius    float64 `json:"temp_c"`
	Fahrenheit float64 `json:"temp_f"`
}

func ValidarCEP(cep string) bool {
	re := regexp.MustCompile(`^\d{5}-?\d{3}$`)
	return re.MatchString(cep)
}

func KelvinBy(c float64) float64 {
	return c + 273
}

var ErrCEPInvalid = errors.New("invalid zipcode")
var ErrCEPNotFound = errors.New("can not find zipcode")
var ErrWeatherFound = errors.New("can not find weather")
