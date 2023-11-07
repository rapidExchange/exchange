package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	RedisUser     string `mapstructure:"REDIS_USER"`
	RedisPassword string `mapstructure:"REDIS_PASS"`
	MysqlUser     string `mapstructure:"MYSQL_USER"`
	MysqlPassword string `mapstructure:"MYSQL_PASSWORD"`
	MysqlDBName   string `mapstructure:"MYSQL_DBNAME"`
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
