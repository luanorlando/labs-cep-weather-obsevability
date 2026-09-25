package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/luanorlando/labs-cep-weather-obsevability.git/internal/entity"
	"github.com/luanorlando/labs-cep-weather-obsevability.git/internal/usecase"
)

type WeatherHandler struct {
	usecase *usecase.FetchWeatherUsecase
}

func NewHandler(u *usecase.FetchWeatherUsecase) *WeatherHandler {
	return &WeatherHandler{
		usecase: u,
	}
}

func (h *WeatherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cep := r.PathValue("cep")

	result, err := h.usecase.Execute(cep)

	if err != nil {
		if errors.Is(err, entity.ErrCEPInvalid) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		if errors.Is(err, entity.ErrCEPNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound) // <-- CORRIGIDO PARA 404
			return
		}

		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(result)
}
