package repository

import (
	"encoding/json"
	"net/http"

	"github.com/luanorlando/labs-cep-weather-obsevability.git/internal/entity"
)

type CEPServiceClient struct {
	httpClient *http.Client
	url        string
}

func NewServiceBClient(url string, client *http.Client) *CEPServiceClient {
	return &CEPServiceClient{
		httpClient: client,
		url:        url,
	}
}

func (c CEPServiceClient) Fetch(cep string) (*entity.Cep, error) {
	// urlAPI := fmt.Sprintf("https://viacep.com.br/ws/%s/json", cep)

	resp, err := c.httpClient.Get(c.url)

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
