package entity

import (
	"errors"
)

type Weather struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_k"`
}

type WeatherFromCity struct {
	Current CurrentWeather `json:"current"`
}

type CurrentWeather struct {
	Celsius    float64 `json:"temp_c"`
	Fahrenheit float64 `json:"temp_f"`
}

func KelvinBy(c float64) float64 {
	return c + 273
}

var ErrWeatherFound = errors.New("can not find weather")
