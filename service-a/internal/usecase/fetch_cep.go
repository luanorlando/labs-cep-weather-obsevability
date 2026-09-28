package usecase

import (
	"github.com/luanorlando/labs-cep-weather-obsevability.git/service-a/internal/repository"
)

type CEPInputDto struct {
	Cep string `json:"cep"`
}

type FetchCepUseCase struct {
	Client *repository.CEPServiceClient
}

func NewFetchCepUseCase(c *repository.CEPServiceClient) *FetchCepUseCase {
	return &FetchCepUseCase{
		Client: c,
	}
}

func (u FetchCepUseCase) Execute(cep string) (string, error) {
	result, err := u.Client.Fetch(cep)
	if err != nil {
		return "", err
	}

	return result.City, nil
}
