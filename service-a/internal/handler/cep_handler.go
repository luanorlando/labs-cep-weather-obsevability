package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/luanorlando/labs-cep-weather-obsevability.git/internal/entity"
)

type CEPDto struct {
	Cep string `json:"cep"`
}

type CepHandler struct {
	serviceBURL string
}

func NewCepHandler() *CepHandler {
	url := os.Getenv("SERVICE_WEATHER_URL")
	if url == "" {
		url = "http://localhost:8082"
	}

	return &CepHandler{
		serviceBURL: url,
	}
}

func (h *CepHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var dto CEPDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, entity.ErrCEPInvalid.Error(), http.StatusUnprocessableEntity)
		return
	}

	if !entity.ValidarCEP(dto.Cep) {
		http.Error(w, entity.ErrCEPInvalid.Error(), http.StatusUnprocessableEntity)
		return

	}
	defer r.Body.Close()

	weatherUrl := fmt.Sprintf("%s/weather?cep=%s", h.serviceBURL, dto.Cep)
	fmt.Printf("URL: %s", weatherUrl)

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

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Erro no serviço B", resp.StatusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("erro ao transmitir dados para o cliente: %v", err)
		return
	}
}
