package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/luanorlando/labs-cep-weather-obsevability.git/config"
	"github.com/luanorlando/labs-cep-weather-obsevability.git/service-a/internal/handler"
	"github.com/luanorlando/labs-cep-weather-obsevability.git/service-a/internal/telemetry"
)

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	collectorURL := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if collectorURL == "" {
		collectorURL = "localhost:4317" // fallback local caso rode fora do Docker
	}

	tp, err := telemetry.InitTracer("service-a", collectorURL)
	if err != nil {
		log.Fatalf("Erro ao inicializar OpenTelemetry: %v", err)
	}

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Erro ao desligar TracerProvider: %v", err)
		}
	}()

	handler := handler.NewCepHandler()

	http.Handle("POST /weather/cep", handler)

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
