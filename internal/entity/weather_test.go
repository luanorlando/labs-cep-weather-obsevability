package entity

import "testing"

func TestValidarCep(t *testing.T) {
	tests := []struct {
		name     string
		inputCEP string
		expected bool
	}{
		{
			name:     "CEP válido padrão",
			inputCEP: "01001000",
			expected: true,
		},
		{
			name:     "CEP Válido com hífen",
			inputCEP: "01001-000",
			expected: true, // Se sua função aceitar/tratar hífens
		},
		{
			name:     "CEP Inválido - Poucos dígitos",
			inputCEP: "123",
			expected: false,
		},
		{
			name:     "CEP Inválido - Letras no meio",
			inputCEP: "0100a000",
			expected: false,
		},
		{
			name:     "CEP Inválido - Vazio",
			inputCEP: "",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(T *testing.T) {
			result := ValidarCEP(tc.inputCEP)
			if result != tc.expected {
				t.Errorf("Para o CEP '%s': esperava %v, mas obteve %v", tc.inputCEP, tc.expected, result)
			}
		})
	}
}

func TestKelvinBy(t *testing.T) {
	tests := []struct {
		name         string
		celsius      float64
		expectedKelv float64
	}{
		{
			name:         "Zero absoluto em Celsius",
			celsius:      -273,
			expectedKelv: 0.0,
		},
		{
			name:         "Ponto de congelamento da água",
			celsius:      0.0,
			expectedKelv: 273,
		},
		{
			name:         "Temperatura ambiente padrão",
			celsius:      25.0,
			expectedKelv: 298,
		},
		{
			name:         "Temperatura do cenário anterior (São Paulo)",
			celsius:      11.1,
			expectedKelv: 284.10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := KelvinBy(tc.celsius)
			// Nota: floats podem ter微as dízimas, se precisar de precisão pode usar uma margem de erro,
			// mas se for soma simples (+ 273.15), a comparação direta funciona bem.
			if result != tc.expectedKelv {
				t.Errorf("Para %.2f°C: esperava %.2fK, mas obteve %.2fK", tc.celsius, tc.expectedKelv, result)
			}
		})
	}
}
