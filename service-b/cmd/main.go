package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/luanorlando/labs-cep-weather-obsevability.git/config"
	"github.com/luanorlando/labs-cep-weather-obsevability.git/service-b/internal/handler"
	"github.com/luanorlando/labs-cep-weather-obsevability.git/service-b/internal/repository"
	"github.com/luanorlando/labs-cep-weather-obsevability.git/service-b/internal/telemetry"
	"github.com/luanorlando/labs-cep-weather-obsevability.git/service-b/internal/usecase"
)

func main() {
	collectorURL := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if collectorURL == "" {
		collectorURL = "localhost:4317"
	}

	tp, err := telemetry.InitTracer("service-b", collectorURL)
	if err != nil {
		log.Fatal("Erro ao inicializar OpenTelemetry no Serviço B: %v", err)
	}

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Erro ao desligar TracerProvider: %v", err)
		}
	}()

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		apiKey = config.WeatherAPIKey
	}

	client := repository.NewWeatherServiceClient(apiKey)
	usecase := usecase.NewFetchWeatherUsecase(client)
	handler := handler.NewHandler(usecase)

	http.Handle("GET /weather", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = config.HTTPPortWeather
	}

	httpPort := fmt.Sprintf(":%s", port)
	log.Printf("Servidor iniciando na porta %s...", httpPort)

	if err := http.ListenAndServe(httpPort, nil); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
