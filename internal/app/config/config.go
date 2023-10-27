package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Logger     Logger     `mapstructure:"logger"      validate:"required"`
	HTTPServer HTTPServer `mapstructure:"httpServer" validate:"required"`
	Postgres   Postgres   `mapstructure:"postgres"   validate:"required"`
	NatsServer NatsServer `mapstructure:"natsServer"   validate:"required"`
}

type (
	Logger struct {
		Level *int8 `mapstructure:"level" validate:"required"`
	}

	Postgres struct {
		Host     string `mapstructure:"host" validate:"required"`
		DBName   string `mapstructure:"dbname" validate:"required"`
		User     string `mapstructure:"user" validate:"required"`
		Password string `mapstructure:"password" validate:"required"`
	}

	HTTPServer struct {
		Address          string `mapstructure:"address" validate:"required"`
		CacheSize        int    `mapstructure:"cacheSize" validate:"required"`
		CacheFillTimeout int64  `mapstructure:"fillTimeout" validate:"required"`
	}
	NatsServer struct {
		ClusterId string `mapstructure:"clusterId" validate:"required"`
		ClientId  string `mapstructure:"clientId" validate:"required"`
		NatsUrl   string `mapstructure:"url" validate:"required"`
	}
)

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")

	defaultLevel := int8(-1)

	viper.SetDefault("logger.level", &defaultLevel)

	viper.SetDefault("postgres.host", "127.0.0.1")
	viper.SetDefault("postgres.dbname", "postgres")
	viper.SetDefault("postgres.user", "postgres")
	viper.SetDefault("postgres.password", "postgres")

	viper.SetDefault("httpServer.address", "0.0.0.0:8080")
	viper.SetDefault("httpServer.cacheSize", "128")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
