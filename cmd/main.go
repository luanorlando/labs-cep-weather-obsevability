package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/luanorlando/labs-cep-weather-obsevability.git/config"
	"github.com/luanorlando/labs-cep-weather-obsevability.git/internal/handler"
	"github.com/luanorlando/labs-cep-weather-obsevability.git/internal/repository"
	"github.com/luanorlando/labs-cep-weather-obsevability.git/internal/usecase"
)

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	client := http.Client{}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		apiKey = config.WeatherAPIKey
	}

	cepRepo := repository.NewViaCEPRepository(&client)
	weatherRepo := repository.NewExternalWeatherRepository(apiKey, &client)
	usecase := usecase.NewFetchWeatherUsecase(cepRepo, weatherRepo)
	cepHandler := handler.NewHandler(usecase)

	http.Handle("GET /weather/{cep}", cepHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = config.HTTPPort
	}

	httpPort := fmt.Sprintf(":%s", port)
	log.Printf("Servidor iniciando na porta %s...", httpPort)
	if err := http.ListenAndServe(httpPort, nil); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
