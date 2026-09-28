package handler

import (
	"encoding/json"
	"net/http"

	"github.com/luanorlando/labs-cep-weather-obsevability.git/internal/entity"
	"github.com/luanorlando/labs-cep-weather-obsevability.git/service-b/internal/usecase"
)

type WeatherHandler struct {
	usecase *usecase.FetchWeather
}

func NewHandler(u *usecase.FetchWeather) WeatherHandler {
	return WeatherHandler{
		usecase: u,
	}
}

func (h WeatherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cepParam := r.URL.Query().Get("cep")
	if cepParam == "" {
		http.Error(w, entity.ErrCEPInvalid.Error(), http.StatusUnprocessableEntity)
		return
	}

	result, err := h.usecase.Execute(cepParam)
	if err != nil {
		http.Error(w, entity.ErrCEPInvalid.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(result)
}
