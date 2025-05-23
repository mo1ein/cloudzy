package config

import (
	_ "github.com/rs/zerolog"
	_ "github.com/rs/zerolog/log"
)

type Config struct {
	//LogLevel log
	Locale   string
	Database Database
	HTTP     HTTP
	GRPC     GRPC
	Tz       string
	BadgeSvc BadgeSvc
}

type Database struct {
	Redis Redis
}

type Redis struct {
	Host     string
	Port     int
	Password string
	Database int
}

type HTTP struct {
	APIHost string
	APIPort int
}

type BadgeSvc struct {
	BaseURL string
}
type GRPC struct {
	Host string
	Port int
}
