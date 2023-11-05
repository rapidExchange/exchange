package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	RedisUser     string `mapstructure:"REDIS_USER"`
	RedisPassword string `mapstructure:"REDIS_PASS"`
}

func LoadConfig(path string) (c Config, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)

	return
}
