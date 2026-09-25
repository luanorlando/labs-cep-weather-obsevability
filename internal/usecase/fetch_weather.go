package usecase

import "github.com/luanorlando/labs-cep-weather-obsevability.git/internal/entity"

type CEPRepository interface {
	Fetch(cep string) (*entity.Cep, error)
}

type WeatherRepository interface {
	FetchBy(city string) (*entity.WeatherFromCity, error)
}

type FetchWeatherUsecase struct {
	cepRepository     CEPRepository
	weatherRepository WeatherRepository
}

func NewFetchWeatherUsecase(cr CEPRepository, wr WeatherRepository) *FetchWeatherUsecase {
	return &FetchWeatherUsecase{
		cepRepository:     cr,
		weatherRepository: wr,
	}
}

func (u FetchWeatherUsecase) Execute(cep string) (*entity.Weather, error) {
	if !entity.ValidarCEP(cep) {
		return nil, entity.ErrCEPInvalid
	}

	cResult, err := u.cepRepository.Fetch(cep)
	if err != nil {
		return nil, err
	}

	wResult, err := u.weatherRepository.FetchBy(cResult.City)
	if err != nil {
		return nil, err
	}

	kelvin := entity.KelvinBy(wResult.Current.Celsius)

	return &entity.Weather{
		Celsius:    wResult.Current.Celsius,
		Fahrenheit: wResult.Current.Fahrenheit,
		Kelvin:     kelvin,
	}, nil
}
