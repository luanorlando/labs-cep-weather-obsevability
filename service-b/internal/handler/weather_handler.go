package handler

import (
	"encoding/json"
	"net/http"

	"github.com/luanorlando/labs-cep-weather-obsevability.git/internal/entity"
	"github.com/luanorlando/labs-cep-weather-obsevability.git/service-b/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

	tracer := otel.Tracer("service-b-handler")
	ctx, span := tracer.Start(ctx, "ExecuteWeatherOrchestration")
	defer span.End()

	cepParam := r.URL.Query().Get("cep")

	if cepParam == "" {
		http.Error(w, entity.ErrCEPInvalid.Error(), http.StatusUnprocessableEntity)
		return
	}

	climaCtx, climaSpan := tracer.Start(ctx, "FetchWeatherExternalAPI")
	result, err := h.usecase.Execute(climaCtx, cepParam)
	climaSpan.End()

	if err != nil {
		http.Error(w, entity.ErrCEPInvalid.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(result)
}
