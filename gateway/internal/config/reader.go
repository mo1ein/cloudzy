package config

import (
	"errors"
	"fmt"
	//log "git.o3social.app/backend/packages/golang/go-logger"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.AllowEmptyEnv(true)

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return nil, fmt.Errorf("reading config: %w", err)
		}
	}
	return &Config{
		Locale: "", //todo: fix this
		Tz:     "", // todo: fix here
		//Database: Database{
		//	MySQL:      nil,
		//	AsynqRedis: nil,
		//},
		HTTP: HTTP{
			APIHost: loadString("HTTP_HOST"),
			APIPort: loadInt("HTTP_PORT"),
		},
		PricingSvc: PricingSvc{
			BaseURL: loadString("PRICING_SVC_BASE_URL"),
		},
		WeatherSvc: WeatherSvc{
			BaseURL: loadString("WEATHER_SVC_BASE_URL"),
		},
	}, nil
}
