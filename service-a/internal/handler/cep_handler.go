package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/luanorlando/labs-cep-weather-obsevability.git/internal/entity"
	"github.com/luanorlando/labs-cep-weather-obsevability.git/service-a/usecase"
)

type CEPDto struct {
	Cep string `json:"cep"`
}

type WeatherOutputDto struct {
	City       string  `json:"city"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_k"`
}

type CepHandler struct {
	usecase     *usecase.FetchCepUseCase
	serviceBURL string
}

func NewCepHandler(u *usecase.) *CepHandler {
	url := os.Getenv("SERVICE_WEATHER_URL")
	if url == "" {
		url = "http://localholst:8082"
	}

	return &CepHandler{
		usecase:     u,
		serviceBURL: url,
	}
}

func (h *CepHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var dto CEPDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}
	// urlAPI := fmt.Sprintf("https://viacep.com.br/ws/%s/json", dto.Cep)
	defer r.Body.Close()

	city, err := h.usecase.Execute(dto.Cep)

	if err != nil {
		if errors.Is(err, entity.ErrCEPInvalid) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		if errors.Is(err, entity.ErrCEPNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	weatherUrl := fmt.Sprintf("%s/weather?city=%s", h.serviceBURL, city)

	weatherReq, err := http.NewRequestWithContext(r.Context(), "GET", weatherUrl, nil)
	if err != nil {
		http.Error(w, "erro ao preparar requisção", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(weatherReq)
	if err != nil {
		http.Error(w, entity.ErrWeatherFound.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var output *WeatherOutputDto

	err = json.NewDecoder(resp.Body).Decode(&output)
	if err != nil {
		http.Error(w, "Erro ao converter dados de temperatura para objeto", http.StatusInternalServerError)
		return
	}

	output.City = city

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(output)
}