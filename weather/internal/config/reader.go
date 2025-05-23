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

	//mysqlConfig := mysql.Config{
	//	Host:         loadString("DATABASE_MYSQL_HOST"),
	//	Port:         loadInt("DATABASE_MYSQL_PORT"),
	//	Username:     loadString("DATABASE_MYSQL_USER"),
	//	Password:     loadString("DATABASE_MYSQL_PASSWORD"),
	//	DatabaseName: loadString("DATABASE_MYSQL_NAME"),
	//	Timezone:     loadString("TZ"),
	//}
	//
	//redisConfig := redis.Config{
	//	Host:     loadString("DATABASE_ASYNCQ_REDIS_HOST"),
	//	Port:     loadString("DATABASE_ASYNCQ_REDIS_PORT"),
	//	Password: loadString("DATABASE_ASYNCQ_REDIS_PASSWORD"),
	//	Database: loadInt("DATABASE_ASYNCQ_REDIS_DATABASE"),
	//}

	return &Config{
		Locale: "", //todo: fix this
		//Database: Database{
		//	MySQL:      nil,
		//	AsynqRedis: nil,
		//},
		GRPC: GRPC{
			Host: loadString("GRPC_HOST"),
			Port: loadInt("GRPC_PORT"),
		},
		Database: Database{
			Redis: Redis{
				Host:     loadString("REDIS_HOST"),
				Port:     loadInt("REDIS_PORT"),
				Password: loadString("REDIS_PASSWORD"),
				Database: loadInt("REDIS_DATABASE"),
			},
		},
		Tz: "", // todo: fix here
		BadgeSvc: BadgeSvc{
			BaseURL: "",
		},
	}, nil
}
