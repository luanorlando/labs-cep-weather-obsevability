package entity

import (
	"errors"
	"regexp"
)

type Cep struct {
	City string `json:"localidade"`
	Erro string `json:"erro,omitempty"`
}

func ValidarCEP(cep string) bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(cep)
}

var ErrCEPInvalid = errors.New("invalid zipcode")
var ErrCEPNotFound = errors.New("can not find zipcode")
