package repository

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/luanorlando/desafio-go-labs-cep-cloudrun.git/internal/entity"
)

type ViaCEPRepository struct {
	httpClient *http.Client
}

func NewViaCEPRepository(client *http.Client) *ViaCEPRepository {
	return &ViaCEPRepository{
		httpClient: client,
	}
}

func (r ViaCEPRepository) Fetch(cep string) (*entity.Cep, error) {
	urlAPI := fmt.Sprintf("https://viacep.com.br/ws/%s/json", cep)

	resp, err := r.httpClient.Get(urlAPI)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info entity.Cep
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	if info.Erro == "true" {
		return nil, entity.ErrCEPNotFound
	}

	return &info, nil
}
