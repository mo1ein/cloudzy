package config

import (
	_ "github.com/rs/zerolog"
	_ "github.com/rs/zerolog/log"
)

type Config struct {
	//LogLevel log
	Locale string
	Tz     string
	//Database Database
	HTTP       HTTP
	GRPC       GRPC
	PricingSvc PricingSvc
	WeatherSvc WeatherSvc
}

type HTTP struct {
	APIHost string
	APIPort int
}

type WeatherSvc struct {
	BaseURL string
}

type PricingSvc struct {
	BaseURL string
}

type GRPC struct {
	Host string
	Port int
}
