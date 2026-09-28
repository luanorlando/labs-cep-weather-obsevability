package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/luanorlando/labs-cep-weather-obsevability.git/config"
	"github.com/luanorlando/labs-cep-weather-obsevability.git/service-b/internal/handler"
	"github.com/luanorlando/labs-cep-weather-obsevability.git/service-b/internal/repository"
	"github.com/luanorlando/labs-cep-weather-obsevability.git/service-b/internal/usecase"
)

func main() {
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
