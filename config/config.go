package config

import (
	"log"

	"github.com/spf13/viper"
)

type Conf struct {
	WeatherAPIKey string `mapstructure:"weather_api_key"`
	HTTPPort      string `mapstructure:"http_port"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg Conf

	viper.SetConfigFile(path + "/.env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	_ = viper.BindEnv("weather_api_key", "WEATHER_API_KEY")
	_ = viper.BindEnv("http_port", "HTTP_PORT")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Aviso: Arquivo .env não encontrado. Buscando variáveis da memória...")
	}

	// O Viper joga as variáveis da memória (injetadas pelo Docker) para dentro da sua struct
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
