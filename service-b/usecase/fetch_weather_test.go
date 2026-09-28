package usecase

import (
	"errors"
	"testing"

	"github.com/luanorlando/labs-cep-weather-obsevability.git/internal/entity"
)

// 1. Definição dos Mocks que satisfazem as interfaces do seu Use Case
type mockCEPRepository struct {
	response *entity.Cep
	err      error
}

func (m *mockCEPRepository) Fetch(cep string) (*entity.Cep, error) {
	return m.response, m.err
}

type mockWeatherRepository struct {
	response *entity.WeatherFromCity
	err      error
}

func (m *mockWeatherRepository) FetchBy(city string) (*entity.WeatherFromCity, error) {
	return m.response, m.err
}

func TestFetchWeatherUsecase_Execute(t *testing.T) {

	type testCase struct {
		name            string
		inputCEP        string
		mockCEPResp     *entity.Cep
		mockCEPErr      error
		mockWeatherResp *entity.WeatherFromCity
		mockWeatherErr  error
		expectErr       error
		expectResult    *entity.Weather
	}

	tests := []testCase{
		{
			name:        "Sucesso - CEP de São Paulo válido retorna temperaturas corretas",
			inputCEP:    "01001000",
			mockCEPResp: &entity.Cep{City: "São Paulo"},
			mockCEPErr:  nil,
			mockWeatherResp: &entity.WeatherFromCity{
				Current: struct {
					Celsius    float64 `json:"temp_c"`
					Fahrenheit float64 `json:"temp_f"`
				}{Celsius: 11.1, Fahrenheit: 52.0},
			},
			mockWeatherErr: nil,
			expectErr:      nil,
			expectResult: &entity.Weather{
				Celsius:    11.1,
				Fahrenheit: 52.0,
				Kelvin:     284.10,
			},
		},
		{
			name:            "Erro - CEP com formato inválido (Tratado pela sua Entidade)",
			inputCEP:        "123",
			mockCEPResp:     nil,
			mockCEPErr:      nil,
			mockWeatherResp: nil,
			mockWeatherErr:  nil,
			expectErr:       entity.ErrCEPInvalid,
			expectResult:    nil,
		},
		{
			name:            "Erro - CEP válido mas não encontrado no banco/API externa",
			inputCEP:        "99999999",
			mockCEPResp:     nil,
			mockCEPErr:      errors.New("cep nao encontrado"), // Simula falha do repositório
			mockWeatherResp: nil,
			mockWeatherErr:  nil,
			expectErr:       errors.New("cep nao encontrado"),
			expectResult:    nil,
		},
		{
			name:            "Erro - Falha ao consultar o serviço de Clima (WeatherRepository)",
			inputCEP:        "01001000",
			mockCEPResp:     &entity.Cep{City: "São Paulo"},
			mockCEPErr:      nil,
			mockWeatherResp: nil,
			mockWeatherErr:  errors.New("api do tempo fora do ar"),
			expectErr:       errors.New("api do tempo fora do ar"),
			expectResult:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			crMock := &mockCEPRepository{response: tc.mockCEPResp, err: tc.mockCEPErr}
			wrMock := &mockWeatherRepository{response: tc.mockWeatherResp, err: tc.mockWeatherErr}

			usecase := NewFetchWeatherUsecase(crMock, wrMock)

			result, err := usecase.Execute(tc.inputCEP)

			if tc.expectErr != nil {
				if err == nil {
					t.Fatalf("Esperava o erro '%v', mas obteve sucesso", tc.expectErr)
				}
				if err.Error() != tc.expectErr.Error() {
					t.Errorf("Esperava o erro '%v', mas obteve '%v'", tc.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Não esperava nenhum erro, mas obteve: %v", err)
			}

			if result.Celsius != tc.expectResult.Celsius {
				t.Errorf("Celsius incorreto: esperava %.1f, obteve %.1f", tc.expectResult.Celsius, result.Celsius)
			}

			if result.Fahrenheit != tc.expectResult.Fahrenheit {
				t.Errorf("Fahrenheit incorreto: esperava %.1f, obteve %.1f", tc.expectResult.Fahrenheit, result.Fahrenheit)
			}

			if result.Kelvin != tc.expectResult.Kelvin {
				t.Errorf("Kelvin incorreto: esperava %.2f, obteve %.2f", tc.expectResult.Kelvin, result.Kelvin)
			}
			_ = result
		})
	}
}
