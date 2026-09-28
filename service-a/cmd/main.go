package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/luanorlando/labs-cep-weather-obsevability.git/config"
	"github.com/luanorlando/labs-cep-weather-obsevability.git/service-a/internal/handler"
)

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

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
